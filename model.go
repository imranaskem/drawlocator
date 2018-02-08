package main

type person struct {
	ID          string `json:"id,omitempty"`
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	PlaceOfWork string `json:"placeofwork,omitempty"`
}
