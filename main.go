/*
 * Alexandria CMDB - Open source common.management database
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
	"github.com/cavaliercoder/alexandria/common"
	"github.com/cavaliercoder/alexandria/controllers"
	"github.com/cavaliercoder/alexandria/database"
	"github.com/cavaliercoder/alexandria/services"

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
			Name:  "c, config",
			Usage: "common.file",
		},
		cli.StringFlag{
			Name:  "answers",
			Usage: "initial common.answer file",
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
		_, err = common.GetConfigFromFile(confFile)
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
		log.Fatal("An answer file was specified but the database is already initialized")
	}

	if !booted {
		if answerFile == "" {
			log.Fatal("Database is not initialized but no answer file was specified.")
		}

		log.Print("Bootstrapping database schema...")
		answers, err := common.LoadAnswers(answerFile)
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

	// Initialize controllers
	controllers := []controllers.Controller{
		new(controllers.CITypeController),
		new(controllers.ConfigController),
		new(controllers.DatabaseController),
		new(controllers.TenantController),
		new(controllers.UserController),
	}

	for _, controller := range controllers {
		m.Group(controller.GetPath(), func(r martini.Router) {
			err = controller.InitRoutes(r)
			if err != nil {
				log.Fatal(err)
			}	
		})
	}

	// Git'er done
	log.Printf("Initialization complete")
	m.Run()
}
