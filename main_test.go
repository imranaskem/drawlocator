package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}

	a.Initialise("TEST", "TEST", "TEST")

	code := m.Run()

	os.Exit(code)
}

func TestGetAll(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "" {
		t.Errorf("Expected content missing")
	}
}

func TestSingleGet(t *testing.T) {
	req := httptest.NewRequest("GET", "/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "" {
		t.Errorf("Expected content missing")
	}
}

func TestSingleUpdate(t *testing.T) {
	updateJSON := json.RawMessage(`{"placeofwork": "Client Office"}`)
	requestBody := bytes.NewBuffer(updateJSON)
	req := httptest.NewRequest("PATCH", "/4", requestBody)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var person person
	_ = json.NewDecoder(response.Body).Decode(&person)

	if person.ID != "4" && person.PlaceOfWork != "Client Office" {
		t.Error("Expected content missing ", person)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
