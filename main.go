package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

//Person is a Draw colleague
type Person struct {
	ID          string `json:"id,omitempty"`
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	PlaceOfWork string `json:"placeofwork,omitempty"`
}

var staff []Person

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", GetAllStaffLocations).Methods("GET")
	router.HandleFunc("/{id}", GetStaffLocation).Methods("GET")
	router.HandleFunc("/{id}", UpdateStaffLocation).Methods("PATCH")

	staff = append(staff, Person{"1", "Kent", "Valentine", "Weston Street"})
	staff = append(staff, Person{"2", "Dean", "Faulkner", "Baker Street"})
	staff = append(staff, Person{"3", "Sian", "Barlow", "Client Office"})
	staff = append(staff, Person{"4", "Imran", "Askem", "Holiday"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}

//GetAllStaffLocations gets locations of all staff
func GetAllStaffLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(staff)
}

//GetStaffLocation gets a single colleague location
func GetStaffLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	for _, item := range staff {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

//UpdateStaffLocation updates a single colleague with a new location
func UpdateStaffLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	for _, item := range staff {
		if item.ID == params["id"] {
			item.PlaceOfWork = person.PlaceOfWork
		}
	}
}
