package main

import (
	"strings"
)

type person struct {
	ID          string `json:"id,omitempty"`
	FirstName   string `json:"firstname,omitempty"`
	LastName    string `json:"lastname,omitempty"`
	PlaceOfWork string `json:"placeofwork,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

func (p *person) generateLocationMessage() string {
	if p.PlaceOfWork == "Weston Street" || p.PlaceOfWork == "Baker Street" {
		return p.FirstName + " is in " + p.PlaceOfWork
	}

	if p.PlaceOfWork == "Holiday" {
		return p.FirstName + " is on " + strings.ToLower(p.PlaceOfWork)
	}

	if p.PlaceOfWork == "Sick" {
		return p.FirstName + " is off " + strings.ToLower(p.PlaceOfWork)
	}

	if p.PlaceOfWork == "Client Office" {
		return p.FirstName + " is at a " + strings.ToLower(p.PlaceOfWork)
	}

	//represents Working from Home
	return p.FirstName + " is " + strings.ToLower(p.PlaceOfWork)

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

type slackResponseMessage struct {
	ResponseType string `json:"response_type"`
	MessageBody  string `json:"text"`
}

func newSlackResponseMessage(msg string) slackResponseMessage {
	s := slackResponseMessage{
		ResponseType: "ephemeral",
		MessageBody:  msg,
	}

	return s
}
