package db

import (
	"fmt"
	"log"
	"strconv"
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

func mkFiscalyear(year int) FiscalYear {
	// the year this document was cretead, it is within a start and end date of a fiscal year
	const dateFmt = "2006-01-02 15:04"
	startDateStr := fmt.Sprintf("%d-01-01 00:00", year)
	endDateStr := fmt.Sprintf("%d-12-31 23:59", year)

	startDate, err := time.Parse(dateFmt, startDateStr)

	if err != nil {
		log.Fatal("error formating fiscalyear date")
	}

	endtDate, err := time.Parse(dateFmt, endDateStr)

	if err != nil {
		log.Fatal("error formating fiscalyear date")
	}

	return FiscalYear{Start: startDate, End: endtDate}
}

func MongoToPgsql(ruc string) {
	// moves data from mongoDB to structured in PostgresSQl Db:
	fmt.Println("Moving data to Postgres database")
	mongodb := GetMongoConnection()

	lieDetails := GetDetailsDataRuc(mongodb, ruc)
	informante := lieDetails.Informante

	// psql connection
	db := GetPGConnection()
	// close connection
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	periodo := 2019
	fiscalYear := mkFiscalyear(periodo)
	fmt.Println(fiscalYear)

	getOrCreateFiscalYear(db, &fiscalYear)

	partyType := PartyType{Type: informante.Clasificacion}
	db.Where("type", partyType.Type).First(&partyType)
	if partyType.ID == 0 {
		db.Create(&partyType)
	}

	juridico := PartyType{Type: "JURIDICO"}
	db.Create(&juridico)

	dv, _ := strconv.Atoi(informante.Dv)
	party := Party{Name: informante.Nombre, TaxPayerId: informante.Ruc,
		DV: dv, PartyType: partyType}

	getOrCreateParty(db, &party)
	var incomes []Income
	for _, jsIncome := range lieDetails.Ingresos {

		customer := Party{
			TaxPayerId: jsIncome.RelacionadoNumeroIdentificacion,
			Name:       jsIncome.RelacionadoNombres,
			PartyType:  partyType,
		}
		getOrCreateParty(db, &customer)

		income := Income{
			Customer:     customer,
			PITIncome:    decimal.NewFromInt(jsIncome.IngresoMontoGravado),
			NotPITIncome: decimal.NewFromInt(jsIncome.IngresoMontoNoGravado),
		}
		income.Documentid = jsIncome.TimbradoDocumento
		incomes = append(incomes, income)
	}

	var expenses []Expense
	for _, jsExpense := range lieDetails.Egresos {

		supplier := Party{
			Name:       jsExpense.RelacionadoNombres,
			TaxPayerId: jsExpense.RelacionadoNumeroIdentificacion,
			PartyType:  juridico,
		}
		getOrCreateParty(db, &supplier)

		expense := Expense{Supplier: supplier}
		expense.TotalAmount = decimal.NewFromInt(jsExpense.EgresoMontoTotal)
		expense.Documentid = jsExpense.TimbradoDocumento
		expenses = append(expenses, expense)

	}
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
