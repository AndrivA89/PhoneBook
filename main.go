package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Структура для контакта
type Contact struct {
	Id          int
	Name        string
	PhoneNumber string
}

// Структура для использования подключения к БД в обработчиках
type Handler struct {
	DB *sqlx.DB
}

// Обработчик главной страницы - вывод всех контактов и телефонов
func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	contacts := []*Contact{}
	rows, _ := h.DB.Query("SELECT `contacts`.`id`, `contacts`.`name`, `phone_number`.`number` FROM `contacts`, `phone_number` WHERE `contacts`.`id` = `phone_number`.`contactsID`;")

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

// Добавление нового контакта
func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	currentContact := &Contact{}
	err := json.NewDecoder(r.Body).Decode(&currentContact)
	errorMsg(err, "Получение JSON - функция Create")
}

func Update(w http.ResponseWriter, r *http.Request) {}

func Delete(w http.ResponseWriter, r *http.Request) {}

// errorMsg - Печать ошибки
func errorMsg(err error, comment string) {
	if err != nil {
		log.Printf("Ошибка %v! Текст ошибки: %v", comment, err)
	}
}

func main() {
	ConnectionStringDB := "root:1234567890@tcp(localhost:3306)/phone_book"
	conn, err := sqlx.Connect("mysql", ConnectionStringDB)
	if err != nil {
		panic(err)
	}

	handlers := &Handler{DB: conn}

	r := mux.NewRouter()

	r.HandleFunc("/contact/new", Create).Methods("POST")
	r.HandleFunc("/contact/{id}", Update).Methods("POST")
	r.HandleFunc("/contact/{id}", Delete).Methods("DELETE")
	r.HandleFunc("/", handlers.MainPage).Methods("GET")

	log.Println("Сервер запущен на порту :88")
	log.Fatal(http.ListenAndServe(":88", r))
}
