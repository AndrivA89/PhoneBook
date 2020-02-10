package handlers

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
var Err error

// Обработчик главной страницы - вывод всех контактов и телефонов
func MainPage(w http.ResponseWriter, r *http.Request) {
	contacts := []*Contact{}
	// Создание запроса для выборки всех контактов
	query := "SELECT `contacts`.`id`, `phone_number`.`id`, `contacts`.`name`, `phone_number`.`number` "
	query += "FROM `contacts`, `phone_number` "
	query += "WHERE `contacts`.`id` = `phone_number`.`contactsID`;"

	rows, err := DB.Query(query)
	ErrorMsg(err, "Отправка запроса выбора данных из таблиц - функция MainPage")

	for rows.Next() {
		currentContact := &Contact{}
		err = rows.Scan(&currentContact.IdContact, &currentContact.IdPhoneNumber, &currentContact.Name, &currentContact.PhoneNumber)
		ErrorMsg(err, "Сканирование строк ответа - функция MainPage")
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

	Err = json.NewDecoder(r.Body).Decode(&currentContact)
	ErrorMsg(Err, "Получение JSON - функция Create")

	// Добавление нового контакта
	res, err := DB.Exec("INSERT INTO `contacts` (name) VALUES(?);", currentContact.Name)
	ErrorMsg(err, "Добавление нового контакта (имя контакта) - функция Create")

	// Получение id нового контакта
	id, err := res.LastInsertId()
	ErrorMsg(err, "Получение индекса последнего контакта - функция Create")

	// Добавление номера телефона нового контакта
	_, err = DB.Exec("INSERT INTO `phone_number` (contactsID, number) VALUES(?, ?);",
		strconv.Itoa(int(id)), currentContact.PhoneNumber)
	ErrorMsg(err, "Добавление нового контакта (ID и номер телефона) - функция Create")
}

// Добавление номера к существующему контакту
func AddNumber(w http.ResponseWriter, r *http.Request) {
	// Берем id элемента из url
	vars := mux.Vars(r)
	idContact, err := strconv.Atoi(vars["idContact"])
	ErrorMsg(err, "Получение ID контакта - функция AddNumber")

	w.Header().Set("Content-Type", "application/json")

	currentContact := &Contact{}

	err = json.NewDecoder(r.Body).Decode(&currentContact)
	ErrorMsg(err, "Получение JSON - функция AddNumber")

	// Добавление номера телефона
	_, err = DB.Exec("INSERT INTO `phone_number` (contactsID, number) VALUES(?, ?);",
		strconv.Itoa(int(idContact)), currentContact.PhoneNumber)
	ErrorMsg(err, "Добавление номера телефона - функция AddNumber")
}

// Обработчик редактирования контакта
func Update(w http.ResponseWriter, r *http.Request) {
	// Берем id элемента из url
	vars := mux.Vars(r)
	idContact, err := strconv.Atoi(vars["idContact"])
	ErrorMsg(err, "Получение ID контакта - функция Update")
	idPhoneNumber, err := strconv.Atoi(vars["idPhoneNumber"])
	ErrorMsg(err, "Получение ID номера телефона - функция Update")

	w.Header().Set("Content-Type", "application/json")

	updateContact := &Contact{}

	err = json.NewDecoder(r.Body).Decode(&updateContact)
	ErrorMsg(err, "Получение JSON - функция Update")

	if idContact != 0 {
		_, err = DB.Exec("UPDATE `contacts` SET `contacts`.`name` = ? WHERE `contacts`.`id` = ?;",
			updateContact.Name,
			strconv.Itoa(idContact))
		ErrorMsg(err, "Отправка запроса на обновление контакта - функция Update")
	}
	if idPhoneNumber != 0 {
		_, err = DB.Exec("UPDATE `phone_number` SET `phone_number`.`number` = ? WHERE `phone_number`.`id` = ?;",
			updateContact.PhoneNumber,
			strconv.Itoa(idPhoneNumber))
		ErrorMsg(err, "Отправка запроса на обновление номера телефона - функция Update")
	}
}

// Обработчик удаления контакта или номера телефона
func Delete(w http.ResponseWriter, r *http.Request) {
	// Берем id элемента из url
	vars := mux.Vars(r)
	idContact, err := strconv.Atoi(vars["idContact"])
	ErrorMsg(err, "Получение ID контакта - функция Delete")
	idPhoneNumber, err := strconv.Atoi(vars["idPhoneNumber"])
	ErrorMsg(err, "Получение ID номера телефона - функция Delete")

	if idContact != 0 { // Если выбран индекс контакта - удалить контакт полностью
		_, err = DB.Exec("DELETE FROM `contacts` where id = ?", idContact)
		ErrorMsg(err, "Удаление контакта полностью - функция Delete")
		_, err = DB.Exec("DELETE FROM `phone_number` where `phone_number`.`contactsID` = ?", idContact)
		ErrorMsg(err, "Удаление контакта полностью (телефон) - функция Delete")
	} else if idPhoneNumber != 0 { // Если выбран только индекс телефона - удалить номер тел.
		_, err = DB.Exec("DELETE FROM `phone_number` where id = ?", idPhoneNumber)
		ErrorMsg(err, "Удаление номера телефона из контакта - функция Delete")
	}
}

// Обработчик поиска контакта по имени или по номеру телефона
func Find(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	findContact := &Contact{}
	contacts := []*Contact{}

	Err = json.NewDecoder(r.Body).Decode(&findContact)
	ErrorMsg(Err, "Получение JSON - функция Find")

	if findContact.Name != "" {
		// Создание запроса для поиска по имени
		query := "SELECT `contacts`.`id`, `phone_number`.`id`, `contacts`.`name`, `phone_number`.`number` "
		query += "FROM `contacts`, `phone_number` "
		query += "WHERE `contacts`.`name` = \"" + findContact.Name + "\" "
		query += "AND `contacts`.`id` = `phone_number`.`contactsID`;"

		rows, err := DB.Query(query)
		ErrorMsg(err, "Отправка запроса поиска по имени - функция Find")

		for rows.Next() {
			currentContact := &Contact{}
			err = rows.Scan(&currentContact.IdContact, &currentContact.IdPhoneNumber, &currentContact.Name, &currentContact.PhoneNumber)
			ErrorMsg(err, "Сканирование строк ответа - функция Find")
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
		ErrorMsg(err, "Отправка запроса поиска по номеру телефона - функция Find")

		for rows.Next() {
			currentContact := &Contact{}
			err = rows.Scan(&currentContact.IdContact, &currentContact.IdPhoneNumber, &currentContact.Name, &currentContact.PhoneNumber)
			ErrorMsg(err, "Сканирование строк ответа - функция Find")
			contacts = append(contacts, currentContact)
		}
		rows.Close()

		json.NewEncoder(w).Encode(contacts)
	}
}

// ErrorMsg - Печать ошибки
func ErrorMsg(err error, comment string) {
	if err != nil {
		log.Printf("Ошибка!!! %v!\n***Текст ошибки:***\n%v", comment, err)
	}
}
