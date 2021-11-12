package handler

import (
	"net/http"
	"text/template"

	"github.com/demianfe/facturaspy/facturaspy/db"
)

type Invoice struct {
	Name     string
	Party    db.Party
	Expenses []db.Expense
}

func InvoiceForm(w http.ResponseWriter, r *http.Request) {
	tmpl := "facturaspy/templates/invoice.html"
	t := template.Must(template.New("invoice.html").ParseFiles(tmpl))

	conn := db.GetPGConnection()
	p := db.Party{TaxPayerId: "2380186"}
	db.GetTaxPayer(conn, &p)

	invoice := Invoice{
		Name:  "expenses",
		Party: p}
	t.Execute(w, &invoice)
}
