package facturaspy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"github.com/kidstuff/mongostore"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	db "github.com/demianfe/facturaspy/facturaspy/db"
	handler "github.com/demianfe/facturaspy/facturaspy/handler"
)

const staticDir = "./facturaspy/static"

func lieRawHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taxpayerId := vars["taxpayerId"]
		fmt.Printf("the ruc is %s\n", taxpayerId)
		mongodb := db.GetMongoConnection()
		lie := db.GetLIEData(mongodb, taxpayerId)
		// Content-type json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(lie)
	})
}

func ledgerHandler() http.Handler {
	/*
		GET /ledger/{taxpayerId}/year/{year}

		Response:
		{
			"incomes": [
			],
			"expenses": [
			]
		}
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taxPayerId := vars["taxpayerId"]
		yearStr := vars["year"]

		var ledger = db.Ledger{}
		year, _ := strconv.Atoi(yearStr)

		db.GetLedger(&ledger, taxPayerId, year)

		// write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ledger)
	})
}

//InitWebServices initializes everything
func InitWebServices() {
	dbsess, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer dbsess.Close()
	// ouath
	store := mongostore.NewMongoStore(dbsess.DB("auth").C("session"), 3600, true,
		[]byte(AuthKey))
	store.MaxAge(MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = IsProd
	// google auth
	goth.UseProviders(
		google.New(GoogleClientID, GoogleClientSecret, GoogleCallback, "email", "profile"),
	)
	gothic.Store = store

	r := mux.NewRouter()
	// template handlers
	r.HandleFunc("/file-upload/", handler.FileUploadForm)
	r.HandleFunc("/ledger/{taxpayerId}", handler.ShowTaxpayer).Methods("GET")
	// API services
	r.Handle("/api/lie/{taxpayerId}", lieRawHandler())
	r.Handle("/api/ledger/{taxpayerId}/year/{year}", ledgerHandler())
	r.HandleFunc("/api/party/{taxpayerId}", handler.PartyHandler).Methods("POST", "PUT")
	r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler).Methods("GET")
	r.HandleFunc("/auth/google/callback", HandleGoogleCallback(store))
	r.Handle("/aranduka/fileupload", handler.FileUploadHandler())
	// r.HandleFunc("/aranduka/fileupload", UploadFile)
	//static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	http.Handle("/", r)

	const port = "8000"
	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("server is listening on port " + port)
	log.Fatal(s.ListenAndServe())
}
