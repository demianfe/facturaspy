package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	db "github.com/demianfe/facturaspy/facturaspy/db"
	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  string
}

// func toMonthInt(m time.Month) string {
// 	return strconv.Itoa(int(m))
// }

func ShowTaxpayer(w http.ResponseWriter, r *http.Request) {
	partyTmpl := "facturaspy/templates/party.html"
	vars := mux.Vars(r)
	taxpayerId := vars["taxpayerId"]

	party := db.Party{TaxPayerId: taxpayerId}
	conn := db.GetPGConnection()
	// get taxpayer data and feed to the template
	db.GetTaxPayer(conn, &party)
	t := template.Must(template.New("party.html").ParseFiles(partyTmpl))
	t.Execute(w, &party)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, msg)
	}
}

func readParty(party *db.Party, r *http.Request) error {

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &party)
	return err
}

func PartyHandler(w http.ResponseWriter, r *http.Request) {
	var party db.Party
	if r.Method == "POST" {
		//create party
		fmt.Println("TODO: create pearty")
	} else if r.Method == "PUT" {
		//update party
		vars := mux.Vars(r)
		tpId := vars["taxpayerId"]

		err := readParty(&party, r)
		if err != nil {
			log.Fatalln(err)
			errorHandler(w, r, http.StatusInternalServerError, "Internar server error")
			return
		}

		if tpId != party.TaxPayerId {
			errorHandler(w, r, http.StatusBadRequest, "Bad request.")
			log.Fatal("url path taxpayerId and json taxpayerId do not match.")
			return
		}

		party.TaxPayerId = tpId
		conn := db.GetPGConnection()
		db.UpdateParty(conn, &party)
	}
}
