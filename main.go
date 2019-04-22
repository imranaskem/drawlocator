package main

import (
	"os"
)

// our main function
func main() {
	dbname := os.Getenv("MONGO_INITDB_DATABASE")
	dburl := os.Getenv("DRAWLOCATOR_DBURL")
	slackToken := os.Getenv("SLACK_SETLOCATION_TOKEN")
	slackWhereIsToken := os.Getenv("SLACK_WHEREIS_TOKEN")
	slackReqToken := os.Getenv("SLACK_OUT_TOKEN")

	a := NewApp(dbname, dburl, slackToken, slackWhereIsToken, slackReqToken)

	a.Run()
}
