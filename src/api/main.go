package main

import (
	"alexandria/api/controllers"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"log"
)

func main() {
	m := martini.Classic()

	// Connect to MongoDB
	mgoSession, err := mgo.Dial("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer mgoSession.Close()

	db := mgoSession.DB("alexandria")

	// Initialize controllers
	_, err = controllers.NewUserController(m, db)
	if err != nil { log.Panic(err) }
	
	_, err = controllers.NewTenantController(m, db)
	if err != nil { log.Panic(err) }

	// Git'er done
	m.Run()
}
