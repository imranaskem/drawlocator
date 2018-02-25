package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

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

	a.Router.GET("/", a.getAllStaffLocations)
	a.Router.OPTIONS("/", a.handleOptions)

	a.Router.GET("/:id", a.getStaffLocation)
	a.Router.PATCH("/:id", a.updateStaffLocation)
	a.Router.OPTIONS("/:id", a.handleOptions)
}

func (a *App) handleOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "OPTIONS, GET, PATCH")
	c.Status(http.StatusOK)
}

func (a *App) getAllStaffLocations(c *gin.Context) {
	var people []person
	_ = a.DB.C("people").Find(bson.M{}).All(&people)

	c.JSON(http.StatusOK, people)
}

func (a *App) getStaffLocation(c *gin.Context) {
	id := c.Param("id")

	var person person
	_ = a.DB.C("people").Find(bson.M{"id": id}).One(&person)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, person)
}

func (a *App) updateStaffLocation(c *gin.Context) {
	id := c.Param("id")

	var personUpdate person
	_ = c.BindJSON(&personUpdate)

	var existingPerson person
	_ = a.DB.C("people").Find(bson.M{"id": id}).One(&existingPerson)

	existingPerson.PlaceOfWork = personUpdate.PlaceOfWork

	_ = a.DB.C("people").Update(bson.M{"id": id}, &existingPerson)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, existingPerson)
}
