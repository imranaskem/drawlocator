package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

//App holds our application
type App struct {
	Router *mux.Router
}

//Initialise acts as our constructor
func (a *App) Initialise(user, pw, dbname string) {
	a.Router = mux.NewRouter()

	staff = getData()

	a.initialiseRoutes()
}

//Run starts our application
func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func (a *App) initialiseRoutes() {
	a.Router.HandleFunc("/", a.getAllStaffLocations).Methods("GET")
	a.Router.HandleFunc("/{id}", a.getStaffLocation).Methods("GET")
	a.Router.HandleFunc("/{id}", a.updateStaffLocation).Methods("PATCH")
}

func (a *App) getAllStaffLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(staff)
}

func (a *App) getStaffLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	for _, item := range staff {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func (a *App) updateStaffLocation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var personUpdate person
	_ = json.NewDecoder(r.Body).Decode(&personUpdate)

	for i, item := range staff {
		if item.ID == params["id"] {
			staff[i].PlaceOfWork = personUpdate.PlaceOfWork
			personUpdate = staff[i]
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(personUpdate)
}
