package main

type person struct {
	ID          string `json:"id,omitempty"`
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	PlaceOfWork string `json:"placeofwork,omitempty"`
}

func getData() (staff []person) {
	staff = append(staff, person{"1", "Kent", "Valentine", "Weston Street"})
	staff = append(staff, person{"2", "Dean", "Faulkner", "Baker Street"})
	staff = append(staff, person{"3", "Sian", "Barlow", "Client Office"})
	staff = append(staff, person{"4", "Imran", "Askem", "Holiday"})
	return
}
