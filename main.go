package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Contact struct {
	Id          int
	Name        string
	PhoneNumber string
}

type Handler struct {
	DB   *sql.DB
	Tmpl *template.Template
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	contacts := []*Contact{}
	rows, err := h.DB.Query("SELECT id, name, phoneNumber FROM contacts")

	//TODO: Обработка ошибки err
	//...
	//...

	for rows.Next() {
		currentContact := &Contact{}
		err = rows.Scan(&currentContact.Id, &currentContact.Name, &currentContact.PhoneNumber)

		//TODO: Обработка ошибки err
		//...
		//...

		contacts = append(contacts, currentContact)
	}
	rows.Close()

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Contacts []*Contact
	}{
		Contacts: contacts,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {}

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

	handlers := &Handler{
		DB:   db,
		Tmpl: template.Must(template.ParseGlob("../templates/*")),
	}

	r := mux.NewRouter()

	r.HandleFunc("/contact/new", handlers.Add).Methods("POST")
	r.HandleFunc("/contact/{id}", handlers.Edit).Methods("GET")
	r.HandleFunc("/contact/{id}", handlers.Update).Methods("POST")
	r.HandleFunc("/contact/{id}", handlers.Delete).Methods("DELETE")
	r.HandleFunc("/", handlers.MainPage).Methods("GET")

	log.Println("Сервер запущен на порту :80")
	log.Fatal(http.ListenAndServe(":80", r))
}
