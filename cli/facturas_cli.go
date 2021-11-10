package main

import (
	"fmt"
	"os"

	db "github.com/demianfe/facturaspy/facturaspy/db"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage.")
		fmt.Println("initdb: initilizes postgres database. WARNING: deletes everything.")
		fmt.Println("topg: lodas postgresql database with from MongoDb.")
		os.Exit(1)
	}
	switch args[0] {
	case "initdb":
		fmt.Println("Initialization ....")
		db.InitPgDb()
	case "topg":
		if len(args) < 2 {
			fmt.Println("Invalid arguments, ruc is not present.")
			os.Exit(1) // some error...
		}
		ruc := args[1]
		fmt.Printf("Moving data to psql for ruc %s.\n", ruc)
		db.MongoToPgsql(ruc)

	case "ledger":
		// print out ladger data
		//l := facturaspy.Ledger{}
		// facturaspy.GetLedger(&l, doc, year)
		// fmt.Println(l.Incomes)
		fmt.Println("TODO")
	default:
		fmt.Printf("Invalid operation: %s\n", args[0])
	}
}
