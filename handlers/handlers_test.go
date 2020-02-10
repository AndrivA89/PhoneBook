package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var testName = "Тестовое имя"         // Для добавления нового Контакта
var testNumber = "+98887776655"       // Для добавления нового Контакта
var testSecondNumber = "+10001110202" // Для добавления второго номера Контакту
var testContactID = 0                 // При отдельном выполнении тестов установить значение
var testNumberID = 0                  // При отдельном выполнении тестов установить значение

func TestCreateContact(t *testing.T) {
	w := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	r := httptest.NewRequest(http.MethodPost, "/contacts/new", buffer)

	testJSON := "{\"Name\": \"" + testName + "\", \"PhoneNumber\": \"" + testNumber + "\"}"
	ConnectionStringDB := "root:1234567890@tcp(localhost:3306)/phone_book"
	DB, _ = sqlx.Connect("mysql", ConnectionStringDB)

	buffer.WriteString(testJSON)
	Create(w, r)

	err := DB.Get(&testContactID, "SELECT `contacts`.`id` FROM `contacts` where `contacts`.`name` = ?", testName)
	if err != nil {
		t.Errorf("Ошибка при проверке на добавление нового контакта")
	}

	number := ""
	err = DB.Get(&number, "SELECT `phone_number`.`number` FROM `phone_number` where `phone_number`.`contactsID` = ?", testContactID)
	if err != nil {
		t.Errorf("Ошибка при проверке на добавление нового контакта")
	}
	if number != testNumber {
		t.Errorf("Добавлен неверный номер телефона")
	}

	DB.Close()
}

func TestAddNumber(t *testing.T) {
	w := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	r := httptest.NewRequest(http.MethodPost, "/contacts/new/", buffer)

	testJSON := "{\"PhoneNumber\": \"" + testSecondNumber + "\"}"
	ConnectionStringDB := "root:1234567890@tcp(localhost:3306)/phone_book"
	DB, _ = sqlx.Connect("mysql", ConnectionStringDB)

	buffer.WriteString(testJSON)
	vars := map[string]string{}
	vars["idContact"] = strconv.Itoa(testContactID)
	vars["idPhoneNumber"] = "0"
	r = mux.SetURLVars(r, vars)

	AddNumber(w, r)

	query := "SELECT `phone_number`.`id`, `phone_number`.`contactsID`, `phone_number`.`number` "
	query += "FROM `phone_number` "
	query += "WHERE `phone_number`.`contactsID` = \"" + strconv.Itoa(testContactID) + "\";"

	rows, err := DB.Query(query)
	if err != nil {
		t.Errorf("Ошибка при проверке на добавление нового телефонного номера")
	}

	contacts := []*Contact{}

	for rows.Next() {
		currentContact := &Contact{}
		err = rows.Scan(&currentContact.IdPhoneNumber, &currentContact.IdContact, &currentContact.PhoneNumber)
		if err != nil {
			t.Errorf("Ошибка при проверке на добавление нового телефонного номера")
		}
		contacts = append(contacts, currentContact)
	}
	rows.Close()

	if contacts[1].PhoneNumber != testSecondNumber {
		t.Errorf("Ошибка добавления нового телефонного номера")
	}
}

func TestDeleteContact(t *testing.T) {
	ConnectionStringDB := "root:1234567890@tcp(localhost:3306)/phone_book"
	DB, _ = sqlx.Connect("mysql", ConnectionStringDB)

	w := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	r := httptest.NewRequest(http.MethodDelete, "/contacts/", buffer)

	vars := map[string]string{}
	vars["idContact"] = strconv.Itoa(testContactID)
	vars["idPhoneNumber"] = strconv.Itoa(testNumberID)
	r = mux.SetURLVars(r, vars)

	Delete(w, r)

	// Создание запроса для выборки всех контактов
	query := "SELECT `contacts`.`id`, `phone_number`.`id`, `contacts`.`name`, `phone_number`.`number` "
	query += "FROM `contacts`, `phone_number` "
	query += "WHERE `contacts`.`id` = `phone_number`.`contactsID`;"

	contacts := []*Contact{}

	rows, err := DB.Query(query)
	if err != nil {
		t.Errorf("Ошибка при проверке на удаление контакта")
	}

	for rows.Next() {
		currentContact := &Contact{}
		err = rows.Scan(&currentContact.IdContact, &currentContact.IdPhoneNumber, &currentContact.Name, &currentContact.PhoneNumber)
		if err != nil {
			t.Errorf("Ошибка при проверке на удаление контакта")
		}
		contacts = append(contacts, currentContact)
	}
	rows.Close()

	for i := 0; i < len(contacts); i++ {
		if contacts[i].IdContact == testContactID {
			t.Errorf("Ошибка при проверке на удаление контакта")
		}
	}

	DB.Close()
}
