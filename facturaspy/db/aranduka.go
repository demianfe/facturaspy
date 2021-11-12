package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// extracts data from aranduka mongodb and persist as helper tables
func MigrateArandukaData() {
	//
	var expenseTypesMap = make(map[string]int)
	expenseTypesMap["Comprobante de Ingreso de Entidades Públicas"] = 16

	pgdb := GetPGConnection()
	mongodb := GetMongoConnection()
	col := mongodb.Database("facturaspy").Collection("lie_details")
	var res []LIEDetalles

	ctx := context.TODO()
	filter := bson.M{}

	cursor, err := col.Find(ctx, filter)

	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
	}
	err = cursor.All(ctx, &res)
	if err != nil {
		fmt.Println("ERROR decoding:", err)
		panic(err)
	}

	for _, lieD := range res {
		for _, eg := range lieD.Egresos {
			aet := ArandukaExpenseType{Text: eg.TipoTexto,
				ExpenseTypeId: int64(expenseTypesMap[eg.TipoTexto])}
			var c int64
			pgdb.Model(&aet).Where(aet).Count(&c)

			if c == 0 {
				fmt.Printf("Creating:  %s. \n", aet.Text)

				pgdb.Create(&aet)
			}
		}
	}
}
