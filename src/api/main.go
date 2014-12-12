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
