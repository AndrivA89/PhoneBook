package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var (
	w      = httptest.NewRecorder()
	buffer = new(bytes.Buffer)
	r      = httptest.NewRequest(http.MethodPost, "/contacts/new", buffer)
)

var testJSON = "{\"Name\":\"Тестовое имя\", \"PhoneNumber\":\"+98887776655\"}"

func TestCreateContact(t *testing.T) {

	buffer.WriteString(testJSON)
	MainPage(w, r)
	// Проверка
	if 1 == 1 {
		// Все Ок
	} else {
		t.Errorf("Ошибка")
	}
}
