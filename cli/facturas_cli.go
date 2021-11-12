package main

import (
	"fmt"
	"os"
	"strconv"

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
		if len(args) < 3 {
			fmt.Println("Invalid arguments, taxpayer id or fiscal year is not present.")
			fmt.Println("Usage topg 1234654 2021")
			os.Exit(1) // some error...
		}
		ruc := args[1]
		fy, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Fiscal year must be an int.")
		}

		fmt.Printf("Moving data to psql for ruc %s for fiscal year %d.\n", ruc, fy)
		db.MongoToPgsql(ruc, fy)
	case "aranduka":
		// move aranduka metadata to helper tables in pg
		db.MigrateArandukaData()
	case "debug":
		// run some arbiriary function
		db.Debug()

	default:
		fmt.Printf("Invalid operation: %s\n", args[0])
	}
}
