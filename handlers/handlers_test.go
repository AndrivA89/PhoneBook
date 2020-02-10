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

	//TODO: Проверить, что номер добавился и теперь у контакта два номера
}

func TestDeleteContact(t *testing.T) {
	ConnectionStringDB := "root:1234567890@tcp(localhost:3306)/phone_book"
	DB, _ = sqlx.Connect("mysql", ConnectionStringDB)

	w := httptest.NewRecorder()
	buffer := new(bytes.Buffer)
	r := httptest.NewRequest(http.MethodDelete, "/contacts/", buffer)

	vars := map[string]string{
		"idContact":     "8",
		"idPhoneNumber": "0",
	}
	r = mux.SetURLVars(r, vars)

	Delete(w, r)

	DB.Close()
}
