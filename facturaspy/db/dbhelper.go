package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DbHost     = "localhost"
	DbPort     = 8432
	DbUser     = "facturas"
	DbPassword = "facturas"
	DbName     = "facturas"
)

func GetPGConnection() *gorm.DB {
	// TODO: use connection pool
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Panic("Connection error")
	}

	return db
}

func InitPgDb() {
	// Drops All tables and regenerates
	// WARNING: you may lose data.

	db := GetPGConnection()
	// close connection
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	fmt.Println("-- Runing automigration --")
	fmt.Println("Dropping all tables")

	db.Migrator().DropTable(
		&FiscalYear{},
		&Invoice{},
		&InvoiceType{},
		&Party{},
		&PartyType{},
		&Income{},
		&IncomeType{},
		&Expense{},
		&ExpenseType{},
		&Ledger{},
		&Stamp{},
		&ArandukaExpenseType{},
	)

	fmt.Println("Recreating all tables")
	err := db.Migrator().AutoMigrate(
		&FiscalYear{},
		&Invoice{},
		&InvoiceType{},
		&Party{},
		&PartyType{},
		&Income{},
		&IncomeType{},
		&Expense{},
		&ExpenseType{},
		&Ledger{},
		&Stamp{},
		&ArandukaExpenseType{},
	)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Inserting basic data")
	var incomeTypes = []IncomeType{}
	for _, text := range IncomeDocTypes {
		it := IncomeType{}
		it.Text = text
		incomeTypes = append(incomeTypes, it)
	}
	db.Create(&incomeTypes)

	var expenseTypes = []ExpenseType{}
	for _, text := range ExpenseDocTypes {
		et := ExpenseType{}
		et.Text = text
		expenseTypes = append(expenseTypes, et)
	}
	db.Create(&expenseTypes)
}

func formatDate(sdate *string) time.Time {
	// var err error
	const dateFmt = "2006-01-02 15:04"

	if len(*sdate) == 10 {
		*sdate = fmt.Sprintf("%s 00:00", *sdate)
	}

	date, err := time.Parse(dateFmt, *sdate)

	if err != nil {
		//TODO: add contextual info
		log.Fatal("error formating date")
		panic(err)
	}

	return date
}

func getOrCreatePartyType(conn *gorm.DB, pt *PartyType) {
	conn.Where("Type", &pt.Type).First(&pt)

	if pt.ID == 0 {
		conn.Create(&pt)
	}

}

func getOrCreateParty(db *gorm.DB, party *Party) {
	// looks into the database for an existing party with its taxpayer id
	// if it does not exists create it
	db.Where("tax_payer_id", &party.TaxPayerId).First(&party)
	if party.ID == 0 {
		fmt.Print("--- Creating a new Party ---\n")
		fmt.Printf("party id: %s, party name: %s\n", party.TaxPayerId, party.Name)
		db.Create(&party)
	}
}

func getOrCreateFiscalYear(db *gorm.DB, fiscalYear *FiscalYear) {
	var count int64
	db.Model(&fiscalYear).
		Where("start", fiscalYear.Start).
		Where("end", fiscalYear.End).
		Count(&count)

	if count == 0 {
		db.Create(&fiscalYear)
	}
}

func getOrCreateStamp(conn *gorm.DB, stamp *Stamp) {
	var count int64
	conn.Model(&stamp).
		Where("party_id = ? and value = ?", stamp.PartyId, stamp.Value).
		Count(&count)

	if count == 0 {
		conn.Create(&stamp)
	}
}

func mkFiscalyear(year int) FiscalYear {
	// the year this document was cretead, it is within a start and end date of a fiscal year

	startDateStr := fmt.Sprintf("%d-01-01 00:00", year)
	endDateStr := fmt.Sprintf("%d-12-31 23:59", year)

	startDate := formatDate(&startDateStr)
	endtDate := formatDate(&endDateStr)
	// endtDate, err := time.Parse(dateFmt, endDateStr)
	// if err != nil {
	// 	log.Fatal("error formating fiscalyear date")
	// }

	return FiscalYear{Start: startDate, End: endtDate}
}

var incomeTypes []IncomeType

func getIncomeType(conn *gorm.DB, it *IncomeType) {
	if len(incomeTypes) == 0 {
		result := conn.Find(&incomeTypes)
		if result.Error != nil {
			log.Fatal("Error retrieving Income types")
			return
		}
	}
	// iterate over the incometypes and return the one that matches
	for _, inct := range incomeTypes {
		if strings.EqualFold(inct.Text, it.Text) {
			*it = inct
			return
		}
	}
}

var expenseTypes []ExpenseType

func getExpenseType(conn *gorm.DB, et *ExpenseType) {
	if len(expenseTypes) == 0 {
		result := conn.Find(&expenseTypes)
		if result.Error != nil {
			log.Fatal("Error retrieving Income types")
			return
		}
	}
	// iterate over the incometypes and return the one that matches
	for _, iet := range expenseTypes {
		if strings.EqualFold(iet.Text, et.Text) {
			*et = iet
			return
		}
	}
	// if we did not find in the look in the aranduka table
	aet := ArandukaExpenseType{Text: et.Text}
	conn.Model(&aet).Where(aet).Find(&aet)

	if aet.ExpenseTypeId != 0 {
		et.ID = aet.ExpenseTypeId
		conn.First(&et, "id = ?", aet.ExpenseTypeId)
	}
}

func createIncomes(db *gorm.DB, incomes *[]Income, lieDetails *LIEDetalles, party *Party, partyType *PartyType) {

	for _, jsIncome := range lieDetails.Ingresos {

		customer := Party{
			TaxPayerId: jsIncome.RelacionadoNumeroIdentificacion,
			Name:       jsIncome.RelacionadoNombres,
			PartyType:  *partyType,
		}
		getOrCreateParty(db, &customer)

		income := Income{
			Customer:     customer,
			PITIncome:    decimal.NewFromInt(jsIncome.IngresoMontoGravado),
			NotPITIncome: decimal.NewFromInt(jsIncome.IngresoMontoNoGravado),
		}

		if jsIncome.TimbradoNumero != "" {
			stamp := Stamp{Value: jsIncome.TimbradoNumero, PartyId: party.ID}
			getOrCreateStamp(db, &stamp)
			income.Stamp = stamp
		}

		income.DocumentId = jsIncome.TimbradoDocumento
		income.Date = (BaseDate)(formatDate(&jsIncome.Fecha))

		// income type
		it := IncomeType{}
		it.Text = jsIncome.TipoTexto
		getIncomeType(db, &it)
		if it.ID != 0 {
			income.IncomeTypeId = it.ID
		}

		*incomes = append(*incomes, income)
	}
}

func createExpenses(db *gorm.DB, expenses *[]Expense, lieDetails *LIEDetalles, party *Party, partyType *PartyType) {

	pt := PartyType{Type: "JURIDICO"}
	getOrCreatePartyType(db, &pt)

	for _, jsExpense := range lieDetails.Egresos {

		provider := Party{
			Name:       jsExpense.RelacionadoNombres,
			TaxPayerId: jsExpense.RelacionadoNumeroIdentificacion,
			PartyType:  pt,
		}
		getOrCreateParty(db, &provider)
		// stamp
		stamp := Stamp{Value: jsExpense.TimbradoNumero, Party: provider}
		getOrCreateStamp(db, &stamp)

		expense := Expense{Provider: provider}
		expense.StampId = stamp.ID
		expense.TotalAmount = decimal.NewFromInt(jsExpense.EgresoMontoTotal)

		expense.DocumentId = jsExpense.TimbradoDocumento
		expense.Date = (BaseDate)(formatDate(&jsExpense.Fecha))

		// expense type
		et := ExpenseType{}
		et.Text = jsExpense.TipoTexto
		getExpenseType(db, &et)
		if et.ID != 0 {
			expense.ExpenseTypeId = et.ID
		}

		*expenses = append(*expenses, expense)
	}
}

func MongoToPgsql(taxpayerId string, fy int) {
	// receives the taxpayer id and the fiscalyear to migrate
	// moves data from mongoDB to structured in PostgresSQl Db:
	fmt.Println("Moving data to Postgres database")
	mongodb := GetMongoConnection()

	lieDetails := GetDetailsDataRuc(mongodb, taxpayerId)
	informante := lieDetails.Informante

	// psql connection
	db := GetPGConnection()
	// close connection
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	fiscalYear := mkFiscalyear(fy)
	fmt.Println(fiscalYear)

	getOrCreateFiscalYear(db, &fiscalYear)

	partyType := PartyType{Type: informante.Clasificacion}
	db.Where("type", partyType.Type).First(&partyType)
	if partyType.ID == 0 {
		db.Create(&partyType)
	}

	dv, _ := strconv.Atoi(informante.Dv)
	party := Party{Name: informante.Nombre, TaxPayerId: informante.Ruc,
		DV: dv, PartyType: partyType}

	getOrCreateParty(db, &party)

	//incomes
	var incomes []Income
	createIncomes(db, &incomes, &lieDetails, &party, &partyType)

	//expenses
	var expenses []Expense
	createExpenses(db, &expenses, &lieDetails, &party, &partyType)

	// create ledger for this party
	ledger := Ledger{
		Party:      party,
		FiscalYear: fiscalYear,
		Incomes:    incomes,
		Expenses:   expenses,
	}
	db.Create(&ledger)
}

func GetLedger(ledger *Ledger, taxpayerId string, year int) {
	// Return all expenses given a taxpayer id and a fiscal year
	db := GetPGConnection()
	// close connection
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// get fiscal year
	fy := mkFiscalyear(year)
	db.First(&fy)

	db.
		Preload("Incomes").
		Preload("Expenses").
		Joins("Party", db.Where("tax_payer_id", taxpayerId)).
		Where("fiscal_year_start", fy.Start).
		Where("fiscal_year_end", fy.End).
		First(&ledger)
}

func GetTaxPayer(conn *gorm.DB, party *Party) {
	// summary data for the taxpayer
	// db := GetPGConnection()
	// close connection
	sqlDB, _ := conn.DB()
	defer sqlDB.Close()
	conn.Where("tax_payer_id", &party.TaxPayerId).First(&party)
}

func UpdateParty(conn *gorm.DB, party *Party) {
	conn.Model(&party).Where("tax_payer_id = ?", party.TaxPayerId).Updates(party)
}

// debug
// call from here arbiriary function or query to test
func Debug() {
	db := GetPGConnection()
	t := ExpenseType{}
	t.Text = "Comprobante de Ingreso de Entidades Públicas"
	fmt.Println(t.Text)
	getExpenseType(db, &t)
	fmt.Println(t.ID)
}
