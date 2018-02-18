package main

import (
	"os"
)

var staff []person

// our main function
func main() {
	a := App{}

	user := os.Getenv("DRAWLOCATOR_USER")
	pw := os.Getenv("DRAWLOCATOR_PW")
	dbname := os.Getenv("DRAWLOCATOR_DBNAME")
	dburl := os.Getenv("DRAWLOCATOR_DBURL")

	a.Initialise(user, pw, dbname, dburl)

	a.Run()
}
