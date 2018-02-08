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

	a.Initialise(user, pw, "drawlocator-db")

	a.Run()
}
