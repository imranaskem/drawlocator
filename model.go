package main

type person struct {
	ID          string `json:"id,omitempty"`
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	PlaceOfWork string `json:"placeofwork,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

type slackUserResponse struct {
	User slackUser `json:"user"`
}

type slackUser struct {
	Profile slackProfile `json:"profile"`
}

type slackProfile struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
