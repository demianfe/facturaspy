package handler

import (
	"fmt"
	"net/http"
	"text/template"
)

func FileUploadForm(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/view/"):]
	p := Page{Title: "LIE Upload"}
	t, _ := template.ParseFiles("facturaspy/templates/file_upload.html")
	t.Execute(w, &p)
}

func FileUploadHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the request comes from the form redirect to details
		// else just return ok
		fmt.Printf("url is %s\n", r.URL)
		taxpayerId, err := UploadFile(w, r)
		if err != nil {
			panic(err)
		}

		fmt.Println(taxpayerId)
		newUrl := fmt.Sprintf("/ledger/%s", taxpayerId)
		http.Redirect(w, r, newUrl, http.StatusSeeOther)
	})
}
