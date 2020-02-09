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
	IdContact     int    `json:"IdContact,omitempty"`
	IdPhoneNumber int    `json:"IdPhoneNumber,omitempty"`
	Name          string `json:"Name"`
	PhoneNumber   string `json:"PhoneNumber"`
}

// Переменная для подключения к БД
var DB *sqlx.DB

// Переменная для получения данных об ошибках
var err error

// Обработчик главной страницы - вывод всех контактов и телефонов
func MainPage(w http.ResponseWriter, r *http.Request) {
	contacts := []*Contact{}
	// Создание запроса для выборки всех контактов
	query := "SELECT `contacts`.`id`, `phone_number`.`id`, `contacts`.`name`, `phone_number`.`number` "
	query += "FROM `contacts`, `phone_number` "
	query += "WHERE `contacts`.`id` = `phone_number`.`contactsID`;"

	rows, err := DB.Query(query)
	errorMsg(err, "Отправка запроса выбора данных из таблиц - функция MainPage")

	for rows.Next() {
		currentContact := &Contact{}
		err = rows.Scan(&currentContact.IdContact, &currentContact.IdPhoneNumber, &currentContact.Name, &currentContact.PhoneNumber)
		errorMsg(err, "Сканирование строк ответа - функция MainPage")
		contacts = append(contacts, currentContact)
	}
	rows.Close()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

// Создание нового контакта
func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	currentContact := &Contact{}

	err = json.NewDecoder(r.Body).Decode(&currentContact)
	errorMsg(err, "Получение JSON - функция Create")

	// Добавление нового контакта
	res, err := DB.Exec("INSERT INTO `contacts` (name) VALUES(?);", currentContact.Name)
	errorMsg(err, "Добавление нового контакта (имя контакта) - функция Create")

	// Получение id нового контакта
	id, err := res.LastInsertId()
	errorMsg(err, "Получение индекса последнего контакта - функция Create")

	// Добавление номера телефона нового контакта
	_, err = DB.Exec("INSERT INTO `phone_number` (contactsID, number) VALUES(?, ?);",
		strconv.Itoa(int(id)), currentContact.PhoneNumber)
	errorMsg(err, "Добавление нового контакта (ID и номер телефона) - функция Create")
}

// Добавление номера к существующему контакту
func AddNumber(w http.ResponseWriter, r *http.Request) {
	// Берем id элемента из url
	vars := mux.Vars(r)
	idContact, err := strconv.Atoi(vars["idContact"])
	errorMsg(err, "Получение ID контакта - функция AddNumber")

	w.Header().Set("Content-Type", "application/json")

	currentContact := &Contact{}

	err = json.NewDecoder(r.Body).Decode(&currentContact)
	errorMsg(err, "Получение JSON - функция AddNumber")

	// Добавление номера телефона
	_, err = DB.Exec("INSERT INTO `phone_number` (contactsID, number) VALUES(?, ?);",
		strconv.Itoa(int(idContact)), currentContact.PhoneNumber)
	errorMsg(err, "Добавление номера телефона - функция AddNumber")
}

// Обработчик редактирования контакта
func Update(w http.ResponseWriter, r *http.Request) {
	// Берем id элемента из url
	vars := mux.Vars(r)
	idContact, err := strconv.Atoi(vars["idContact"])
	errorMsg(err, "Получение ID контакта - функция Update")
	idPhoneNumber, err := strconv.Atoi(vars["idPhoneNumber"])
	errorMsg(err, "Получение ID номера телефона - функция Update")

	w.Header().Set("Content-Type", "application/json")

	updateContact := &Contact{}

	err = json.NewDecoder(r.Body).Decode(&updateContact)
	errorMsg(err, "Получение JSON - функция Update")

	if idContact != 0 {
		_, err = DB.Exec("UPDATE `contacts` SET `contacts`.`name` = ? WHERE `contacts`.`id` = ?;",
			updateContact.Name,
			strconv.Itoa(idContact))
		errorMsg(err, "Отправка запроса на обновление контакта - функция Update")
	}
	if idPhoneNumber != 0 {
		_, err = DB.Exec("UPDATE `phone_number` SET `phone_number`.`number` = ? WHERE `phone_number`.`id` = ?;",
			updateContact.PhoneNumber,
			strconv.Itoa(idPhoneNumber))
		errorMsg(err, "Отправка запроса на обновление номера телефона - функция Update")
	}
}

// Обработчик удаления контакта
func Delete(w http.ResponseWriter, r *http.Request) {
	// Берем id элемента из url
	vars := mux.Vars(r)
	idContact, err := strconv.Atoi(vars["idContact"])
	errorMsg(err, "Получение ID контакта - функция Delete")
	idPhoneNumber, err := strconv.Atoi(vars["idPhoneNumber"])
	errorMsg(err, "Получение ID номера телефона - функция Delete")

	if idContact != 0 { // Если выбран индекс контакта - удалить контакт полностью
		_, err = DB.Exec("DELETE FROM `contacts` where id = ?", idContact)
		errorMsg(err, "Удаление контакта полностью - функция Delete")
	} else if idPhoneNumber != 0 { // Если выбран только индекс телефона - удалить номер тел.
		_, err = DB.Exec("DELETE FROM `phone_number` where id = ?", idPhoneNumber)
		errorMsg(err, "Удаление номера телефона из контакта - функция Delete")
	}
}

// Обработчик поиска контакта по имени или по номеру телефона
func Find(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	findContact := &Contact{}
	contacts := []*Contact{}

	err = json.NewDecoder(r.Body).Decode(&findContact)
	errorMsg(err, "Получение JSON - функция Find")

	if findContact.Name != "" {
		// Создание запроса для поиска по имени
		query := "SELECT `contacts`.`id`, `phone_number`.`id`, `contacts`.`name`, `phone_number`.`number` "
		query += "FROM `contacts`, `phone_number` "
		query += "WHERE `contacts`.`name` = \"" + findContact.Name + "\" "
		query += "AND `contacts`.`id` = `phone_number`.`contactsID`;"

		rows, err := DB.Query(query)
		errorMsg(err, "Отправка запроса поиска по имени - функция Find")

		for rows.Next() {
			currentContact := &Contact{}
			err = rows.Scan(&currentContact.IdContact, &currentContact.IdPhoneNumber, &currentContact.Name, &currentContact.PhoneNumber)
			errorMsg(err, "Сканирование строк ответа - функция Find")
			contacts = append(contacts, currentContact)
		}
		rows.Close()

		json.NewEncoder(w).Encode(contacts)
	} else if findContact.PhoneNumber != "" {
		// Создание запроса для поиска по номеру телефона
		query := "SELECT `contacts`.`id`, `phone_number`.`id`, `contacts`.`name`, `phone_number`.`number` "
		query += "FROM `contacts`, `phone_number` "
		query += "WHERE `phone_number`.`number` = \"" + findContact.PhoneNumber + "\" "
		query += "AND `contacts`.`id` = `phone_number`.`contactsID`;"

		rows, err := DB.Query(query)
		errorMsg(err, "Отправка запроса поиска по номеру телефона - функция Find")

		for rows.Next() {
			currentContact := &Contact{}
			err = rows.Scan(&currentContact.IdContact, &currentContact.IdPhoneNumber, &currentContact.Name, &currentContact.PhoneNumber)
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
	r.HandleFunc("/contacts/addNumber/{idContact}", AddNumber).Methods("POST")
	r.HandleFunc("/contacts/{idContact}/{idPhoneNumber}", Update).Methods("POST")
	r.HandleFunc("/contacts/{idContact}/{idPhoneNumber}", Delete).Methods("DELETE")
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
