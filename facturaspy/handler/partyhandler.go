package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	db "github.com/demianfe/facturaspy/facturaspy/db"
	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  string
}

func toMonthInt(m time.Month) string {
	return strconv.Itoa(int(m))
}

func ShowTaxpayer(w http.ResponseWriter, r *http.Request) {
	partyTmpl := "facturaspy/templates/party.html"
	vars := mux.Vars(r)
	taxpayerId := vars["taxpayerId"]
	// get taxpayer data and feed to the template
	party := db.Party{TaxPayerId: taxpayerId}
	db.GetTaxPayer(&party)

	fmap := template.FuncMap{
		"toMonthInt": toMonthInt,
	}

	t := template.Must(template.New("party.html").Funcs(fmap).ParseFiles(partyTmpl))
	t.Execute(w, &party)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, msg string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, msg)
	}
}

func PartyHandler(w http.ResponseWriter, r *http.Request) {
	var party db.Party
	if r.Method == "POST" {
		//create party
		fmt.Println("TODO: create pearty")
	} else if r.Method == "PUT" {
		//update party
		b, err1 := io.ReadAll(r.Body)
		if err1 != nil {
			log.Fatalln(err1)
		}

		fmt.Println(string(b))
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&party)

		defer r.Body.Close()
		if err != nil {
			fmt.Println(err)
			errorHandler(w, r, http.StatusInternalServerError, "Internar server error")
			panic(err)
			// return
		}
		fmt.Println(party)

	}
}
