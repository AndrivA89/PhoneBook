package main

import (
	"log"
	"net/http"

	handlers "github.com/PhoneBook/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Структура для контакта
type Contact struct {
	IdContact     int    `json:"IdContact,omitempty"`
	IdPhoneNumber int    `json:"IdPhoneNumber,omitempty"`
	Name          string `json:"Name"`
	PhoneNumber   string `json:"PhoneNumber"`
}

// Переменная для подключения к БД
var DB *sqlx.DB

// Переменная для получения данных об ошибках
var err error

func main() {
	ConnectionStringDB := "root:1234567890@tcp(localhost:3306)/phone_book"
	DB, err = sqlx.Connect("mysql", ConnectionStringDB)
	errorMsg(err, "Подключение к БД - функция main")

	r := mux.NewRouter()

	r.HandleFunc("/contacts/new", handlers.Create).Methods("POST")
	r.HandleFunc("/contacts/new/{idContact}", AddNumber).Methods("POST")
	r.HandleFunc("/contacts/{idContact}/{idPhoneNumber}", Update).Methods("POST")
	r.HandleFunc("/contacts/{idContact}/{idPhoneNumber}", Delete).Methods("DELETE")
	r.HandleFunc("/contacts/find", Find).Methods("GET")
	r.HandleFunc("/contacts/", MainPage).Methods("GET")
	r.HandleFunc("/", MainPage).Methods("GET")

	log.Println("Сервер запущен на порту :88")
	log.Fatal(http.ListenAndServe(":88", r))
}

// errorMsg - Печать ошибки
func errorMsg(err error, comment string) {
	if err != nil {
		log.Printf("Ошибка!!! %v!\n***Текст ошибки:***\n%v", comment, err)
	}
}
