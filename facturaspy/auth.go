package facturaspy

import (
	"context"
	"fmt"
	"net/http"

	"log"

	db "github.com/demianfe/facturaspy/facturaspy/db"
	"github.com/gorilla/sessions"
	"github.com/kidstuff/mongostore"
	"github.com/markbates/goth/gothic"
)

const (
	//AuthKey used in gothig
	AuthKey = "facturaspy-session-key" // Replace with your SESSION_SECRET or similar
	//MaxAge auth key duration
	MaxAge = 86400 * 30 // 30 days
	//IsProd change to true for production
	IsProd = false // Set to true when serving over https
	//GoogleClientID for facutraspy
	GoogleClientID = ""
	//GoogleClientSecret for facutraspy
	GoogleClientSecret = ""
	//GoogleCallback for facutraspy
	GoogleCallback = "http://localhost:8000/auth/google/callback"
)

//var lock = &sync.Mutex{} // locks goroutines
//var store *mongostore.MongoStore

// func getStore() *mongostore.MongoStore {
// 	if store == nil {
// 		lock.Lock()
// 		defer lock.Unlock()
// 		if store == nil {
// 			fmt.Println("Creting Single Instance Now")
// 			dbsess, err := mgo.Dial("localhost")
// 			if err != nil {
// 				panic(err)
// 			}
// 			defer dbsess.Close()
// 			store = mongostore.NewMongoStore(dbsess.DB("auth").C("session"), 3600, true,
// 				[]byte(key))
// 			store.MaxAge(maxAge)
// 			store.Options.Path = "/"
// 			store.Options.HttpOnly = true // HttpOnly should always be enabled
// 			store.Options.Secure = isProd
// 		} else {
// 			fmt.Println("Single Instance already created-1")
// 		}
// 	} else {
// 		fmt.Println("Single Instance already created-2")
// 	}
// 	return store
// }

//HandleGoogleCallback handles authentication through google
func HandleGoogleCallback(store *mongostore.MongoStore) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Handling google callback!")
		// Fetch new store.
		//store = getStore()
		// Get a session.
		session, err := store.Get(req, AuthKey)

		if err != nil {
			log.Println(err.Error())
		}

		// Add a value.
		session.Values["foo"] = "bar"

		// Save.
		if err = sessions.Save(req, w); err != nil {
			log.Printf("Error saving session: %v", err)
		}

		fmt.Fprintln(w, "ok")

		user, err := gothic.CompleteUserAuth(w, req)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Println(user)
		mongodb := db.GetMongoConnection()

		usersCol := mongodb.Database("auth").Collection("users")
		insertResult, err := usersCol.InsertOne(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted post with ID:", insertResult.InsertedID)
	})
}

//InitGoth initalizes goth service
// func InitGoth() {
// 	goth.UseProviders(
// 		google.New(googleClientID, googleClientSecret, googleCallback, "email", "profile"),
// 	)
// 	gothic.Store = getStore()

// }

//InitGoogleAuth initiates google login
// func InitGoogleAuth() http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 		gothic.BeginAuthHandler(w, req)
// 	})
// }

// 	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
// 		gothic.BeginAuthHandler(res, req)
// 	})

// 	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
// 		t, _ := template.ParseFiles("templates/index.html")
// 		t.Execute(res, false)
// 	})
// 	// log.Println("listening on localhost:8000")
// 	// log.Fatal(http.ListenAndServe(":8000", p))
// }
