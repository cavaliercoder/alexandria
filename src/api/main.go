/*
 * Alexandria CMDB - Open source configuration management database
 * Copyright (C) 2014  Ryan Armstrong <ryan@cavaliercoder.com>
 * 
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * package controllers
 */
package main

import (
	"alexandria/api/controllers"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"log"
        "alexandria/api/application"
)

func main() {
	// Initialize Martini
        m := martini.Classic()

	// Initialize MongoDB
	mgoSession, err := mgo.Dial("localhost")
	if err != nil {
		log.Panic(err)
	}
	defer mgoSession.Close()

	db := mgoSession.DB("alexandria")
        
        // Initialize application context
        app := application.AppContext{m, db}

	// Initialize controllers
	configController := new(controllers.ConfigController)
	err = configController.Init(&app)
	if err != nil { log.Fatal(err) }
	
	userController := new(controllers.UserController)
        err = userController.Init(&app)
	if err != nil { log.Fatal(err) }
	
	tenantController := new(controllers.TenantController)
        err = tenantController.Init(&app)
	if err != nil { log.Fatal(err) }

	// Git'er done
	m.Run()
}
