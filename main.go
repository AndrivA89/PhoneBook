package main

import (
	"log"
	"net/http"

	h "github.com/PhoneBook/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {
	ConnectionStringDB := "root:1234567890@tcp(localhost:3306)/phone_book"
	h.DB, h.Err = sqlx.Connect("mysql", ConnectionStringDB)
	h.ErrorMsg(h.Err, "Подключение к БД - функция main")

	r := mux.NewRouter()

	r.HandleFunc("/contacts/new", h.Create).Methods("POST")
	r.HandleFunc("/contacts/new/{idContact}", h.AddNumber).Methods("POST")
	r.HandleFunc("/contacts/{idContact}/{idPhoneNumber}", h.Update).Methods("POST")
	r.HandleFunc("/contacts/{idContact}/{idPhoneNumber}", h.Delete).Methods("DELETE")
	r.HandleFunc("/contacts/find", h.Find).Methods("GET")
	r.HandleFunc("/contacts/", h.MainPage).Methods("GET")
	r.HandleFunc("/", h.MainPage).Methods("GET")

	log.Println("Сервер запущен на порту :88")
	log.Fatal(http.ListenAndServe(":88", r))
}
