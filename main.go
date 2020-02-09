package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Contact struct {
	Id          int
	Name        string
	PhoneNumber string
}

var DB *sql.DB

func MainPage(w http.ResponseWriter, r *http.Request) {
	contacts := []*Contact{}
	rows, _ := DB.Query("SELECT id, name, phoneNumber FROM contacts")

	//TODO: Обработка ошибки err (_)
	//...
	//...

	for rows.Next() {
		currentContact := &Contact{}
		_ = rows.Scan(&currentContact.Id, &currentContact.Name, &currentContact.PhoneNumber)

		//TODO: Обработка ошибки err
		//...
		//...

		contacts = append(contacts, currentContact)
	}
	rows.Close()
}

func Add(w http.ResponseWriter, r *http.Request) {}

func Edit(w http.ResponseWriter, r *http.Request) {}

func Update(w http.ResponseWriter, r *http.Request) {}

func Delete(w http.ResponseWriter, r *http.Request) {}

func main() {
	dsn := "root@tcp(localhost)" // Логин-пароль к БД
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()

	r.HandleFunc("/contact/new", Add).Methods("POST")
	r.HandleFunc("/contact/{id}", Edit).Methods("GET")
	r.HandleFunc("/contact/{id}", Update).Methods("POST")
	r.HandleFunc("/contact/{id}", Delete).Methods("DELETE")
	r.HandleFunc("/", MainPage).Methods("GET")

	log.Println("Сервер запущен на порту :80")
	log.Fatal(http.ListenAndServe(":80", r))
}
