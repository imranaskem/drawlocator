package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const slackGetURL = "https://slack.com/api/users.info?token="

//App holds our application
type App struct {
	Router                *gin.Engine
	DB                    *mgo.Database
	SlackSetLocationToken string
	SlackWhereIsToken     string
	SlackRequestToken     string
}

//NewApp acts as our constructor
func NewApp(dbname, dburl, slackSetLocationToken, slackWhereIsToken, slackReqToken string) *App {
	a := App{}

	a.Router = gin.Default()

	a.SlackSetLocationToken = slackSetLocationToken

	a.SlackWhereIsToken = slackWhereIsToken

	a.SlackRequestToken = slackReqToken

	a.initialiseRoutes()

	dialled := false
	count := 1

	for dialled == false {
		fmt.Printf("Connecting to database, attempt %v\n", count)
		s, err := mgo.Dial(dburl)

		if err != nil {
			fmt.Println("Error: " + err.Error())
			fmt.Println("Retrying...")
		} else {
			fmt.Println("Connected!")
			dialled = true
			a.DB = s.DB(dbname)
		}

		time.Sleep(2 * time.Second)

		if count > 10 {
			fmt.Printf("Tried to connect %v times. Exiting program", count)
			os.Exit(1)
		}
		count++
	}

	return &a
}

//Run starts our application
func (a *App) Run() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	a.Router.Run(":" + port)
	fmt.Println("App ready for use")
}

func (a *App) initialiseRoutes() {
	a.Router.GET("/locations", a.getLocations)

	a.Router.GET("/staff", a.getAllStaffLocations)
	a.Router.OPTIONS("/staff", a.handleOptions)

	a.Router.GET("/staff/:id", a.getStaffLocation)
	a.Router.PATCH("/staff/:id", a.updateStaffLocation)
	a.Router.OPTIONS("/staff/:id", a.handleOptions)

	a.Router.GET("/websocket", a.websocketHandler)

	a.Router.POST("/slack", a.handleSlackRequest)
}

func (a *App) handleSlackRequest(c *gin.Context) {
	com := c.Request.FormValue("command")

	switch {
	case com == "/setlocation":
		a.updateLocationFromSlack(c)

	case com == "/whereis":
		a.getLocationFromSlack(c)
	}
}

func (a *App) getLocations(c *gin.Context) {
	locations := []string{"Weston Street", "Baker Street", "Holiday", "Sick", "Client Office", "Working from Home"}

	c.JSON(http.StatusOK, locations)
}

func (a *App) getLocationFromSlack(c *gin.Context) {
	token := c.Request.FormValue("token")

	if token != a.SlackWhereIsToken {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	msg := c.Request.FormValue("text")

	uid, _ := getUserID(msg)

	slackPerson := a.getPersonFromSlackAPI(uid)

	slackMessage := newSlackResponseMessage(slackPerson.generateLocationMessage())

	c.JSON(http.StatusOK, slackMessage)
}

func (a *App) updateLocationFromSlack(c *gin.Context) {
	token := c.Request.FormValue("token")

	if token != a.SlackSetLocationToken {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	msg := c.Request.FormValue("text")
	newLocation, err := standardisePlace(msg)

	if err != nil {
		srm := newSlackResponseMessage("Invalid location")
		c.AbortWithStatusJSON(http.StatusBadRequest, srm)
		return
	}

	userid := c.Request.FormValue("user_id")
	existingPerson := a.getPersonFromSlackAPI(userid)

	existingPerson.PlaceOfWork = newLocation
	a.updatePerson(existingPerson)

	srm := newSlackResponseMessage("Location updated to " + newLocation)

	c.JSON(http.StatusOK, srm)
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
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		if r.Header.Get("Origin") == "app" {
			return true
		}

		return false
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

func (a *App) getAllPeople() []person {
	var people []person
	_ = a.DB.C("people").Find(bson.M{}).All(&people)

	sort.Slice(people, func(i, j int) bool { return people[i].FirstName < people[j].FirstName })

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

func (a *App) findPersonByName(first, last string) person {
	var existingPerson person
	_ = a.DB.C("people").Find(bson.M{"firstname": first, "lastname": last}).One(&existingPerson)

	return existingPerson
}

func (a *App) updatePerson(p person) error {
	return a.DB.C("people").Update(bson.M{"id": p.ID}, &p)
}

func (a *App) getPersonFromSlackAPI(userid string) person {
	resp, _ := http.Get(slackGetURL + a.SlackRequestToken + "&user=" + userid)

	var slackResponse slackUserResponse
	json.NewDecoder(resp.Body).Decode(&slackResponse)

	return a.findPersonByName(slackResponse.User.Profile.FirstName, slackResponse.User.Profile.LastName)
}
