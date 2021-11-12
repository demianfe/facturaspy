package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Model loosely based in
// https://www.set.gov.py/portal/PARAGUAY-SET/detail?folder-id=repository:collaboration:/sites/PARAGUAY-SET/categories/SET/biblioteca-virtual/registros-de-libro-compra-venta-ingreso-egreso/registro-en-planilla-electronica-rg-55-2020&content-id=/repository/collaboration/sites/PARAGUAY-SET/documents/biblioteca/biblioteca-virtual/2020/Modelo%20de%20Libro%20ingreso%20y%20egreso%20para%20quienes%20sean%20solo%20contribuyentes%20del%20%20IRP%20(Prestaci%C3%B3n%20de%20servicios%20en%20relaci%C3%B3n%20de%20dependencia).xlsx
// https://www.set.gov.py/portal/PARAGUAY-SET/detail?folder-id=repository:collaboration:/sites/PARAGUAY-SET/categories/SET/biblioteca-virtual/registros-de-libro-compra-venta-ingreso-egreso/registro-en-planilla-electronica-rg-55-2020&content-id=/repository/collaboration/sites/PARAGUAY-SET/documents/biblioteca/biblioteca-virtual/2020/Modelo%20de%20Libro%20ventas,%20ingresos,%20compras,%20egresos%20para%20contribuyentes%20que%20tengan%20solo%20IVA%20o%20IVA%20y%20Rentas..xlsx

type BaseDate time.Time

func (d BaseDate) Year() int {
	return time.Time(d).Year()
}

func (d BaseDate) Day() int {
	return time.Time(d).Day()
}

func (d BaseDate) Month() int {
	return int(time.Time(d).Month())
}

func (date *BaseDate) UnmarshalJSON(b []byte) error {
	// extract string and transform it to time
	const dateFmt = "2006-01-02 15:04"

	dateStr := strings.Trim(string(b), `"`)
	dateStr = fmt.Sprintf("%s 00:00", dateStr)

	d, err := time.Parse(dateFmt, dateStr)
	if err != nil {
		return err
	}
	*date = (BaseDate)(d)
	return nil
}

type FiscalYear struct {
	Start time.Time `gorm:"primaryKey"`
	End   time.Time `gorm:"primaryKey"`
}

type BaseModel struct {
	ID        int64      `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" gorm:"index"`
}

type InvoiceType struct {
	BaseModel
	Type string
}

type Invoice struct {
	BaseModel
	IssueDate       BaseDate
	InvoiceType     InvoiceType
	TaxPayerId      int64
	InvoiceTypeId   int64
	FiscalYear      FiscalYear `gorm:"ForeignKey:FiscalYearStart,FiscalYearEnd;References:Start,End"`
	FiscalYearStart time.Time
	FiscalYearEnd   time.Time
}

type Stamp struct {
	BaseModel
	Value   string
	Party   Party
	PartyId int64
}

// add party document type
type Party struct {
	BaseModel
	TaxPayerId  string `gorm:"index:idx_taxpayerid,unique" json:"taxpayerId"`
	DV          int
	Name        string    `json:"name"`
	BirthDate   BaseDate  `json:"birthDate"`
	StartDate   time.Time `json:"startDate"`
	PartyType   PartyType
	PartyTypeId int64 `json:"-"`
	Stamps      []Stamp
}

type PartyType struct {
	BaseModel
	Type string `gorm:"index:idx_partytype,unique"`
}

type DocumentType struct {
	BaseModel
	TextId string
	Text   string
}

type User struct {
	gorm.Model
}

type IncomeType struct {
	DocumentType
}

type Document struct {
	gorm.Model
	DocumentId  string
	Date        BaseDate
	TotalAmount decimal.Decimal `json:"totalAmount" sql:"type:decimal(15,2);"`
	TotalExempt decimal.Decimal `json:"totalExepmt" sql:"type:decimal(15,2);"`
	TotalVat5   decimal.Decimal `json:"totalVat5" sql:"type:decimal(15,2);"`
	TotalVat10  decimal.Decimal `json:"totalVat10" sql:"type:decimal(15,2);"`
	TotalNovat  decimal.Decimal `json:"totalNoVat" sql:"type:decimal(15,2);"`
}

type Income struct {
	Document
	PITIncome         decimal.Decimal `sql:"type:decimal(15,2);"`
	NotPITIncome      decimal.Decimal `sql:"type:decimal(15,2);"`
	Customer          Party
	CustomerId        int64
	LedgerId          int64
	Ledger            Ledger
	StampId           int64
	Stamp             Stamp
	DocumentTypeId    int64
	TransactionTypeId int64
	IncomeTypeId      int64
	IncomeType        IncomeType
}

type ExpenseType struct {
	DocumentType
}

type ExpenseSubType struct {
	BaseModel
	TextId string
	Text   string
}

type Expense struct {
	Document
	PITDeductible     decimal.Decimal `sql:"type:decimal(10,2);"` // Personal Income Tax deductible
	Provider          Party
	ProviderId        int64 // party id
	LedgerId          int64
	Ledger            Ledger
	StampId           int64
	Stamp             Stamp
	TransactionTypeId int64
	ExpenseType       ExpenseType
	ExpenseTypeId     int64
}

type Ledger struct {
	BaseModel
	Party           Party `json:"party"`
	PartyId         int64
	OwnerId         int64      // user associated to this ledger
	Incomes         []Income   `json:"incomes"`
	Expenses        []Expense  `json:"expenses"`
	FiscalYear      FiscalYear `json:"fiscalYear" gorm:"ForeignKey:FiscalYearStart,FiscalYearEnd;References:Start,End"`
	FiscalYearStart time.Time  `json:"fiscalYearStart"`
	FiscalYearEnd   time.Time  `json:"fiscalYearEnd"`
}

// helper table
type ArandukaExpenseType struct {
	BaseModel
	Text          string
	ExpenseTypeId int64
}
