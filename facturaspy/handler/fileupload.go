package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	db "github.com/demianfe/facturaspy/facturaspy/db"
	filehelper "github.com/demianfe/facturaspy/facturaspy/filehelper"
	helper "github.com/demianfe/facturaspy/facturaspy/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// upload file using curl:
// curl -X POST -L -F "metadata={name : 'LIE_2019_95da08ce_xxxxxxx_yyy.zip'};type=application/json;charset=UTF-8" -F "file=@LIE_2019_95da08ce_xxxxxxx_yyy.zip;type=application/zip" http://localhost:8000/aranduka/fileupload

func readAndSaveLIE(filename string) string {
	// reads a file from the file system and dumps it to a mongodb object
	// header
	var lie db.LIE

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("File reading error", err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &lie)

	mongoDB := db.GetMongoConnection()
	c := mongoDB.Database(db.MongoDBName).Collection("lie")

	fmt.Printf("Saving LIE Header data for ruc %s\n", lie.Informante.Ruc)
	filterLie := bson.D{{Key: "informante.ruc", Value: lie.Informante.Ruc}}
	updateLie := bson.D{{Key: "$set", Value: lie}}

	opts := options.Update().SetUpsert(true)
	lieRes, _ := c.UpdateOne(context.TODO(), filterLie, updateLie, opts)
	fmt.Println(lieRes)
	return lie.Informante.Ruc
}

func readAndSaveDetails(filename string) {
	// reads a file from the file system and dumps it to a mongodb object
	// details
	var det db.LIEDetalles

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("File reading error", err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &det)

	mongoDB := db.GetMongoConnection()
	c := mongoDB.Database(db.MongoDBName).Collection("lie_details")
	fmt.Println("Saving Lie Details for ruc " + det.Informante.Ruc)

	filterLie := bson.D{{Key: "informante.ruc", Value: det.Informante.Ruc}}
	updateLie := bson.D{{Key: "$set", Value: det}}

	opts := options.Update().SetUpsert(true)
	lieRes, _ := c.UpdateOne(context.TODO(), filterLie, updateLie, opts)
	fmt.Println(lieRes)

}

//UploadFile for file uploading
func UploadFile(w http.ResponseWriter, r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20) // limit max input length
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fmt.Printf("Uploading: %s\n", header.Filename)
	var filePath = filehelper.UploadDir + header.Filename
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.Copy(f, file)

	// unzip file
	dest := strings.Split(header.Filename, ".")
	fileNames, err := filehelper.Unzip(filePath, filehelper.UploadDir+dest[0])

	// check that the json we need is present
	lieFn := filepath.Join(filepath.Clean(filehelper.UploadDir), dest[0], dest[0]+".json")
	detailsFn := filepath.Join(filepath.Clean(filehelper.UploadDir), dest[0], dest[0]+"-detalle.json")

	if helper.Contains(fileNames, lieFn) && helper.Contains(fileNames, detailsFn) {
		ruc := readAndSaveLIE(lieFn)
		readAndSaveDetails(detailsFn)
		db.MongoToPgsql(ruc)
		return ruc, nil
	}
	return "", err
}
