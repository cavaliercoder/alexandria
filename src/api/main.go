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
	"alexandria/api/configuration"
	"alexandria/api/controllers"
	"alexandria/api/database"
	"alexandria/api/services"

	"github.com/codegangsta/cli"
	"github.com/go-martini/martini"

	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Alexandria CMDB API Server Daemon"
	app.Usage = "api"
	app.Version = "1.0.0"
	app.Author = "Ryan Armstrong"
	app.Email = "ryan@cavaliercoder.com"

	// Global args
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "c, config",
			Usage: "configuration file",
		},
		cli.StringFlag{
			Name:  "answers",
			Usage: "initial configuration answer file",
		},
	}

	app.Action = serve
	app.Run(os.Args)

}

func serve(context *cli.Context) {
	var err error
	
	// Load configuration
	confFile := context.GlobalString("config")
	if confFile != "" {
		_, err = configuration.GetConfigFromFile(confFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	// Establish db connection
	log.Printf("Initializing database connection...")
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Check is db schema is initialized
	log.Printf("Checking database schema...")
	booted, err := db.IsBootStrapped()
	if err != nil {
		log.Fatal(err)
	}

	// Build db schema if required
	answerFile := context.GlobalString("answers")
	if booted && answerFile != "" {
		log.Fatal("An answer file was specified by the database is already initialized")	
	}
	
	if !booted {
		if answerFile == "" {
			log.Fatal("Database is not initialized but no answer file was specified.")
		}

		log.Print("Bootstrapping database schema...")
		answers, err := configuration.LoadAnswers(answerFile)
		if err != nil {
			log.Fatal(err)
		}

		err = db.BootStrap(answers)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize Martini
	log.Printf("Starting HTTP handlers...")
	m := martini.Classic()
	m.Use(services.ApiContextService())
	m.Use(services.ApiKeyValidation())

	// Initialize controllers
	initRoutes(m.Router)

	// Git'er done
	log.Printf("Initialization complete")
	m.Run()
}

func initRoutes(router martini.Router) {
	var err error
	configController := new(controllers.ConfigController)
	err = configController.Init(router)
	if err != nil {
		log.Fatal(err)
	}

	userController := new(controllers.UserController)
	err = userController.Init(router)
	if err != nil {
		log.Fatal(err)
	}

	tenantController := new(controllers.TenantController)
	err = tenantController.Init(router)
	if err != nil {
		log.Fatal(err)
	}
}
