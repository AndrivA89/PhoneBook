package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Структура для контакта
type Contact struct {
	Id          int    `json:"Id,omitempty"`
	Name        string `json:"Name"`
	PhoneNumber string `json:"PhoneNumber"`
}

// Переменная для подключения к БД
var DB *sqlx.DB

// Переменная для получения данных об ошибках
var err error

// Обработчик главной страницы - вывод всех контактов и телефонов
func MainPage(w http.ResponseWriter, r *http.Request) {
	contacts := []*Contact{}
	// Создание запроса для выборки всех контактов
	query := "SELECT `contacts`.`id`, `contacts`.`name`, `phone_number`.`number` "
	query += "FROM `contacts`, `phone_number` "
	query += "WHERE `contacts`.`id` = `phone_number`.`contactsID`;"

	rows, err := DB.Query(query)
	errorMsg(err, "Отправка запроса выбора данных из таблиц - функция MainPage")

	for rows.Next() {
		currentContact := &Contact{}
		err = rows.Scan(&currentContact.Id, &currentContact.Name, &currentContact.PhoneNumber)
		errorMsg(err, "Сканирование строк ответа - функция MainPage")
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

	err = json.NewDecoder(r.Body).Decode(&currentContact)
	errorMsg(err, "Получение JSON - функция Create")

	// Добавление нового контакта
	res, err := DB.Exec("INSERT INTO `contacts` (name) VALUES(\"" + currentContact.Name + "\");")
	errorMsg(err, "Добавление нового контакта (имя контакта) - функция Create")

	// Получение id нового контакта
	id, err := res.LastInsertId()
	errorMsg(err, "Получение индекса последнего контакта - функция Create")

	// Добавление номера телефона нового контакта
	query := "INSERT INTO `phone_number` (contactsID, number) VALUES(\"" +
		strconv.Itoa(int(id)) + "\", \"" +
		currentContact.PhoneNumber + "\");"
	_, err = DB.Exec(query)
	errorMsg(err, "Добавление нового контакта (ID и номер телефона) - функция Create")
}

func Update(w http.ResponseWriter, r *http.Request) {}

func Delete(w http.ResponseWriter, r *http.Request) {}

// Обработчик поиска контакта по имени или по номеру телефона
func Find(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	findContact := &Contact{}
	contacts := []*Contact{}

	err = json.NewDecoder(r.Body).Decode(&findContact)
	errorMsg(err, "Получение JSON - функция Find")

	if findContact.Name != "" {
		// Создание запроса для поиска по имени
		query := "SELECT `contacts`.`id`, `contacts`.`name`, `phone_number`.`number` "
		query += "FROM `contacts`, `phone_number` "
		query += "WHERE `contacts`.`name` = \"" + findContact.Name + "\" "
		query += "AND `contacts`.`id` = `phone_number`.`contactsID`;"

		rows, err := DB.Query(query)
		errorMsg(err, "Отправка запроса поиска по имени - функция Find")

		for rows.Next() {
			currentContact := &Contact{}
			err = rows.Scan(&currentContact.Id, &currentContact.Name, &currentContact.PhoneNumber)
			errorMsg(err, "Сканирование строк ответа - функция Find")
			contacts = append(contacts, currentContact)
		}
		rows.Close()

		json.NewEncoder(w).Encode(contacts)
	} else if findContact.PhoneNumber != "" {
		// Создание запроса для поиска по номеру телефона
		query := "SELECT `contacts`.`id`, `contacts`.`name`, `phone_number`.`number` "
		query += "FROM `contacts`, `phone_number` "
		query += "WHERE `phone_number`.`number` = \"" + findContact.PhoneNumber + "\" "
		query += "AND `contacts`.`id` = `phone_number`.`contactsID`;"

		rows, err := DB.Query(query)
		errorMsg(err, "Отправка запроса поиска по номеру телефона - функция Find")

		for rows.Next() {
			currentContact := &Contact{}
			err = rows.Scan(&currentContact.Id, &currentContact.Name, &currentContact.PhoneNumber)
			errorMsg(err, "Сканирование строк ответа - функция Find")
			contacts = append(contacts, currentContact)
		}
		rows.Close()

		json.NewEncoder(w).Encode(contacts)
	}
}

func main() {
	ConnectionStringDB := "root:1234567890@tcp(localhost:3306)/phone_book"
	DB, err = sqlx.Connect("mysql", ConnectionStringDB)
	errorMsg(err, "Подключение к БД - функция main")

	r := mux.NewRouter()

	r.HandleFunc("/contacts/new", Create).Methods("POST")
	r.HandleFunc("/contacts/{id}", Update).Methods("POST")
	r.HandleFunc("/contacts/{id}", Delete).Methods("DELETE")
	r.HandleFunc("/contacts/find", Find).Methods("GET")
	r.HandleFunc("contacts/", MainPage).Methods("GET")
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
