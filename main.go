package main

import (
	"os"
)

var staff []person

// our main function
func main() {
	user := os.Getenv("DRAWLOCATOR_USER")
	pw := os.Getenv("DRAWLOCATOR_PW")
	dbname := os.Getenv("DRAWLOCATOR_DBNAME")
	dburl := os.Getenv("DRAWLOCATOR_DBURL")
	slackToken := os.Getenv("SLACK_TOKEN")
	slackReqToken := os.Getenv("SLACK_OUT_TOKEN")

	a := NewApp(user, pw, dbname, dburl, slackToken, slackReqToken)

	a.Run()
}
