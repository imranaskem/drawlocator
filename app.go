package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//App holds our application
type App struct {
	Router *gin.Engine
	DB     *mgo.Database
}

//Initialise acts as our constructor
func (a *App) Initialise(user, pw, dbname, dburl string) {
	a.Router = gin.Default()

	a.initialiseRoutes()

	s, _ := mgo.Dial(dburl)

	a.DB = s.DB(dbname)

	a.DB.Login(user, pw)
}

//Run starts our application
func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	a.Router.Run(":" + port)
}

func (a *App) initialiseRoutes() {
	a.Router.GET("/staff", a.getAllStaffLocations)
	a.Router.OPTIONS("/staff", a.handleOptions)

	a.Router.POST("/sms", a.smsHandler)

	a.Router.GET("/staff/:id", a.getStaffLocation)
	a.Router.PATCH("/staff/:id", a.updateStaffLocation)
	a.Router.OPTIONS("/staff/:id", a.handleOptions)

	a.Router.GET("/websocket", a.websocketHandler)
}

func (a *App) smsHandler(c *gin.Context) {
	msg, _ := c.GetPostForm("Body")
	from, _ := c.GetPostForm("From")

	place, err := standardisePlace(msg)

	if err != nil {
		c.Header("Access-Control-Allow-Origin", "*")
		c.String(http.StatusOK, "Invalid location: "+place)
		return
	}

	pers := a.findPersonbyPhoneNumber(from)

	pers.PlaceOfWork = place

	a.updatePerson(pers)

	c.Header("Access-Control-Allow-Origin", "*")
	c.String(http.StatusOK, "Location updated to "+place)
}

func (a *App) handleOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "OPTIONS, GET, PATCH")
	c.Status(http.StatusOK)
}

func (a *App) getAllStaffLocations(c *gin.Context) {
	people := a.getAllPeople()

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, people)
}

func (a *App) getStaffLocation(c *gin.Context) {
	id := c.Param("id")

	person := a.findPersonbyID(id)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, person)
}

func (a *App) updateStaffLocation(c *gin.Context) {
	id := c.Param("id")

	var personUpdate person
	_ = c.BindJSON(&personUpdate)

	existingPerson := a.findPersonbyID(id)

	existingPerson.PlaceOfWork = personUpdate.PlaceOfWork

	_ = a.DB.C("people").Update(bson.M{"id": id}, &existingPerson)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, existingPerson)
}

func (a *App) websocketHandler(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     originCheck,
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	oPeople := a.getAllPeople()
	conn.WriteJSON(oPeople)

	for {
		time.Sleep(1 * time.Second)
		nPeople := a.getAllPeople()

		if !comparePeople(oPeople, nPeople) {
			conn.WriteJSON(nPeople)
			oPeople = nPeople
		}
	}
}

func comparePeople(a, b []person) bool {
	return reflect.DeepEqual(a, b)
}

func originCheck(r *http.Request) bool {
	/* ori := r.Header.Get("Origin")

	if ori == "http://websocket.local" {
		return true
	} */

	return true
}

func (a *App) getAllPeople() []person {
	var people []person
	_ = a.DB.C("people").Find(bson.M{}).All(&people)

	return people
}

func (a *App) findPersonbyID(id string) person {
	var existingPerson person
	_ = a.DB.C("people").Find(bson.M{"id": id}).One(&existingPerson)

	return existingPerson
}

func (a *App) findPersonbyPhoneNumber(number string) person {
	var existingPerson person
	_ = a.DB.C("people").Find(bson.M{"phone": number}).One(&existingPerson)

	return existingPerson
}

func (a *App) updatePerson(p person) error {
	return a.DB.C("people").Update(bson.M{"id": p.ID}, &p)
}

func standardisePlace(place string) (string, error) {
	place = strings.ToLower(place)

	switch {
	case strings.Contains(place, "baker"):
		return "Baker Street", nil

	case strings.Contains(place, "sick"):
		return "Sick", nil

	case strings.Contains(place, "weston"):
		return "Weston Street", nil

	case strings.Contains(place, "holiday"):
		return "Holiday", nil

	case strings.Contains(place, "client"):
		return "Client Office", nil

	case strings.Contains(place, "home"):
		return "Working from Home", nil
	}

	return place, errors.New("Invalid place")
}
