package main

import (
	"os"
)

var staff []person

// our main function
func main() {
	a := App{}

	user := os.Getenv("USER")
	pw := os.Getenv("PW")
	dbname := os.Getenv("DBNAME")
	dburl := os.Getenv("DBURL")

	a.Initialise(user, pw, dbname, dburl)

	a.Run()
}
