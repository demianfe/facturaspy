package filehelper

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"

	db "github.com/demianfe/facturaspy/facturaspy/db"
	helper "github.com/demianfe/facturaspy/facturaspy/helper"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

const (
	//UploadDir is where the uploaded zip files recide
	UploadDir = "./files/uploaded/"
	//GeneratedDir is where temporarily created files are sotred
	//this files are then compressed to a zip file
	GeneratedDir = "./files/out/"
)

//GenZIPFile generates zip file to be exported
func GenZIPFile(ruc string) {

	// Read data form db
	// gen filename:
	// <periodo>_<8 random chars>_<ruc>_<obligacion>.xls
	// ej 952 libros irp
	// <uuid>-detalle.json
	// <uuid>-egresos.xls
	// <uuid>-inglresos.xls
	// <uuid>.xml
	// <uuid>.json
	// Finally compress the file using zip
	// encoding ISO-8859-14

	// Ingresos headers:
	// Tipo de Documento	Tipo de Documento (Texto)	Tipo de Ingreso	Tipo de Ingreso (Texto)	Fecha	Mes	Tipo de Identificación	Número de Identificación	Nombres y Apellidos o Razón Social	Número de Timbrado	Número de Documento	Condición de la Venta	Monto Gravado	Monto No Gravado	Monto Total	Número de Cuenta	Razón Social del Banco / Financiera / Cooperativa	Otro Tipo de Documento	Número de Documento
	// Columns with empty values are filled with tab "\t"

	//TODO mkdir

	// ingresosHeader contains the
	// titles that goes into the xls file
	// ingresosJSONKey contains the keys of the
	// the index of the first array maps to the key
	// in the second array which in turn is the key in the json item

	charSet := "0123456789abcdedfghijklmnopqrst"
	var randStr strings.Builder
	for i := 0; i < 8; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		randStr.WriteString(string(randomChar))
	}
	// 952 libros_irp
	var baseFileName = GeneratedDir + "LIE_" + "periodo_" + randStr.String() + "_" + ruc + "_952"

	// TODO: use Go rutines
	fileNamesMap := map[string]string{
		"xls_ingresos": baseFileName + "-ingresos.xls",
		"xls_egresos":  baseFileName + "-egresos.xls",
		"json_details": baseFileName + "-detalle.json",
		"xml":          baseFileName + ".xml",
		"json":         baseFileName + ".json",
	}

	enc := charmap.ISO8859_14.NewEncoder() // latin encoder

	mongodb := db.GetMongoConnection()
	lie := db.GetLIEData(mongodb, ruc)
	lieDet := db.GetDetailsDataRuc(mongodb, ruc)

	//TODO: call this function with using go rutines
	mkXML(fileNamesMap["xml"], &lie)
	mkJSONSummary(fileNamesMap["json"], &lie)
	mkJSONDetails(fileNamesMap["json_details"], &lieDet)
	mkIngresosXLS(enc, fileNamesMap["xls_ingresos"], &lieDet)
	mkEgresosXLS(enc, fileNamesMap["xls_egresos"], &lieDet)

	fileNames := make([]string, 0, len(fileNamesMap))
	for _, v := range fileNamesMap {
		fileNames = append(fileNames, v)
	}

	ZipFiles(baseFileName+".zip", fileNames)
	deleteFiles(fileNames)
}

func deleteFiles(files []string) {
	for _, fn := range files {
		err := os.Remove(fn)
		helper.CheckError(err)
	}
}

//ZipFiles to download
func ZipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = strings.ReplaceAll(filename, GeneratedDir, "")

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func mkJSONDetails(baseFileName string, lieDet *db.LIEDetalles) {
	data, merr := json.MarshalIndent(lieDet, "", " ")
	helper.CheckError(merr)
	_ = ioutil.WriteFile(baseFileName, data, 0644)
}

func mkJSONSummary(fileName string, lie *db.LIE) {
	data, merr := json.MarshalIndent(lie, "", " ")
	helper.CheckError(merr)
	_ = ioutil.WriteFile(fileName, data, 0644)
}

func mkXML(fileName string, lie *db.LIE) {
	xmlData := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`
	x, _ := xml.MarshalIndent(lie, "", "    ")
	xs := string(x)
	xs = strings.ReplaceAll(xs, "LIE>", "resumen>")
	xmlData = xmlData + xs

	f, err := os.Create(fileName)
	helper.CheckError(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	//write headers
	_, werr := w.WriteString(xmlData)
	helper.CheckError(werr)
	w.Flush()
}

func mkIngresosXLS(enc *encoding.Encoder, fileName string, lieDet *db.LIEDetalles) {
	// TODO: handle missing headers
	// array index is related to struct field order
	ingresoHeader := []string{
		"Tipo de Documento",
		"Tipo de Documento (Texto)",
		"Tipo de Ingreso",
		"Tipo de Ingreso (Texto)",
		"Fecha",
		"Mes",
		"Tipo de Identificación",
		"Número de Identificación",
		"Nombres y Apellidos o Razón Social",
		"Número de Timbrado",
		"Número de Documento",
		"Condición de la Venta",
		"Monto Gravado",
		"Monto No Gravado",
		"Monto Total",
		"Número de Cuenta",
		"Razón Social del Banco / Financiera / Cooperativa",
		"Otro Tipo de Documento",
		"Número de Documento"}

	var errH error
	h := strings.Join(ingresoHeader, "\t") + "\n"
	h, errH = enc.String(h)
	helper.CheckError(errH)
	ingresoXLSFile, err := os.Create(fileName)
	helper.CheckError(err)
	defer ingresoXLSFile.Close()
	xlsWriter := bufio.NewWriter(ingresoXLSFile)
	//write headers
	_, werr := xlsWriter.WriteString(h)
	helper.CheckError(werr)

	// TODO: use reflection once we know the structure of the json
	for _, item := range lieDet.Ingresos {
		// ingreso xls
		xlsRow := mkIngresoXLSRow(enc, &item)
		_, err := xlsWriter.WriteString(xlsRow)
		helper.CheckError(err)
	}
	xlsWriter.Flush()
	fmt.Println("Written file " + fileName)
}

func mkIngresoXLSRow(enc *encoding.Encoder, item *db.IngresoDetalle) string {
	ptrStrct := reflect.ValueOf(item)
	strct := ptrStrct.Elem()
	// we know it is a struct but we check it anyways
	var row string
	if strct.Kind() == reflect.Struct {
		var line []string
		for i := 0; i < strct.NumField(); i++ {
			var val string = "\t"
			f := strct.Field(i)
			fn := strct.Type().Field(i).Name

			// ignored fileds: fields present only in json
			if helper.Contains([]string{"ID", "Ruc", "Periodo"}, fn) {
				break
			}
			if f.Kind() == reflect.Int {
				val = strconv.FormatInt(int64(f.Int()), 10)
			} else if f.Kind() == reflect.String {
				val = f.String()
			} else {
				fmt.Println("Unhandled key ", fn, "of kind ", f.Kind())
			}
			line = append(line, val)
		}
		for i := len(line); i < strct.NumField()-1; i++ {
			line = append(line, "\t")
		}
		r, encErr := enc.String(strings.Join(line, "\t") + "\n")
		row = r
		helper.CheckError(encErr)
	}
	return row
}

func mkEgresosXLS(enc *encoding.Encoder, fileName string, lieDet *db.LIEDetalles) {
	//TODO: handle missing headers

	egresoHeader := []string{
		"Tipo de Documento",
		"Tipo de Documento (Texto)",
		"Fecha",
		"Mes",
		"Tipo de Identificación",
		"Número de Identificación",
		"Nombres y Apellidos o Razón Social",
		"Número de Timbrado",
		"Número de Documento",
		"Condición de la Venta",
		"Monto Total",
		"Número de Cuenta",
		"Razón Social del Banco / Financiera / Cooperativa",
		"Otro Tipo de Documento",
		"Número de Documento",
		"Número de Despacho",
		"Período de la Cuenta",
		"Identificador del Empleador",
		"Tipo de Egreso",
		"Tipo de Egreso (Texto)",
		"Clasificación de Egreso",
		"Clasificación de Egreso (Texto)",
		"Número de Identificación Del Empleador"}

	var errH error
	h := strings.Join(egresoHeader, "\t") + "\n"
	h, errH = enc.String(h)
	helper.CheckError(errH)

	xlsFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer xlsFile.Close()

	writer := bufio.NewWriter(xlsFile)
	//write headers
	_, werr := writer.WriteString(h)
	helper.CheckError(werr)

	for _, e := range lieDet.Egresos {
		var line []string

		strct := reflect.ValueOf(&e).Elem()
		for i := 0; i < strct.NumField(); i++ {
			tag := strct.Type().Field(i).Tag.Get("xls")

			if helper.Contains(egresoHeader, tag) {
				var val string = "\t"
				f := strct.Field(i)
				fn := strct.Type().Field(i).Name
				if helper.Contains([]string{"ID", "Ruc", "Periodo"}, fn) {
					continue
				}
				switch f.Kind() {
				case reflect.Int:
					val = strconv.FormatInt(int64(f.Int()), 10)

				case reflect.String:
					val = f.String()
				case reflect.Float32, reflect.Float64:
					fmt.Println("Float value " + f.String())
				}
				line = append(line, val)
			}
		}

		for i := len(line); i < strct.NumField()-1; i++ {
			line = append(line, "\t")
		}
		row, encErr := enc.String(strings.Join(line, "\t") + "\n")
		helper.CheckError(encErr)
		writer.WriteString(row)
	}

	writer.Flush()

}
