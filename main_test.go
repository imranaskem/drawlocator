package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a *App

func TestMain(m *testing.M) {
	dbname := os.Getenv("MONGO_INITDB_DATABASE")
	dburl := os.Getenv("DRAWLOCATOR_DBURL")
	slackToken := os.Getenv("SLACK_SETLOCATION_TOKEN")
	slackWhereIsToken := os.Getenv("SLACK_WHEREIS_TOKEN")
	slackReqToken := os.Getenv("SLACK_OUT_TOKEN")

	a = NewApp(dbname, dburl, slackToken, slackWhereIsToken, slackReqToken)

	code := m.Run()

	os.Exit(code)
}

func TestGetAll(t *testing.T) {
	req := httptest.NewRequest("GET", "/staff", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "" {
		t.Errorf("Expected content missing")
	}
}

func TestSingleGet(t *testing.T) {
	req := httptest.NewRequest("GET", "/staff/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "" {
		t.Errorf("Expected content missing")
	}
}

func TestSingleUpdate(t *testing.T) {
	updateJSON := json.RawMessage(`{"placeofwork": "Client Office"}`)
	requestBody := bytes.NewBuffer(updateJSON)
	req := httptest.NewRequest("PATCH", "/staff/4", requestBody)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var person person
	_ = json.NewDecoder(response.Body).Decode(&person)

	if person.ID != "4" && person.PlaceOfWork != "Client Office" {
		t.Errorf("Expected content missing %v\n", person)
	}
}

func TestComparePeople(t *testing.T) {
	peopleA := createPeople()
	peopleB := createPeople()

	same := comparePeople(peopleA, peopleB)

	if !same {
		t.Error("Contents not equal!")
	}
}

func TestRegex(t *testing.T) {
	s, _ := getUserID("<@U4EMTUT36|imran>")

	if s != "U4EMTUT36" {
		t.Errorf("Regex not working correctly, generated %v\n", s)
	}
}

func createPeople() []person {
	var people []person

	for i := 0; i < 4; i++ {
		person := person{
			ID:          string(i),
			FirstName:   "Test",
			LastName:    "Number " + string(i),
			PlaceOfWork: "Weston Street",
			Phone:       "+447525944042",
		}

		people = append(people, person)
	}

	return people
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
