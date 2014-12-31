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
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"log"
	"os"
	"os/user"
	"regexp"
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

	app.Action = func(context *cli.Context) {
		var err error

		// Load configuration
		confFile := context.GlobalString("config")
		if confFile != "" {
			_, err = GetConfigFromFile(confFile)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Check is db schema is initialized
		log.Printf("Checking database schema...")
		booted, err := IsBootStrapped()
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
			answers, err := LoadAnswers(answerFile)
			if err != nil {
				log.Fatal(err)
			}

			err = BootStrap(answers)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Start web server
		Serve()

		// Git'er done
		log.Printf("Initialization complete")
	}
	app.Run(os.Args)

}

func GetServer() *negroni.Negroni {
	// Init Mux routes
	rootRouter := mux.NewRouter()
	router := rootRouter.PathPrefix(ApiV1Prefix).Subrouter()
	router.HandleFunc("/info", GetApiInfo).Methods("GET")

	// User routes
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users", AddUser).Methods("POST")
	router.HandleFunc("/users/current", GetCurrentUser).Methods("GET")
	router.HandleFunc("/users/{email}", GetUserByEmail).Methods("GET")
	router.HandleFunc("/users/{email}", DeleteUserByEmail).Methods("DELETE")
	router.HandleFunc("/users/{email}/password", SetUserPassword).Methods("PATCH")

	// Tenant routes
	router.HandleFunc("/tenants", GetTenants).Methods("GET")
	router.HandleFunc("/tenants", AddTenant).Methods("POST")
	router.HandleFunc("/tenants/current", GetCurrentTenant).Methods("GET")
	router.HandleFunc("/tenants/{code}", GetTenantByCode).Methods("GET")
	router.HandleFunc("/tenants/{code}", DeleteTenantByCode).Methods("DELETE")

	// CMDB routes
	router.HandleFunc("/cmdbs", GetCmdbs).Methods("GET")
	router.HandleFunc("/cmdbs", AddCmdb).Methods("POST")
	router.HandleFunc("/cmdbs/{name}", GetCmdbByName).Methods("GET")
	router.HandleFunc("/cmdbs/{name}", DeleteCmdbByName).Methods("DELETE")

	// CI Type routes
	router.HandleFunc("/cmdbs/{cmdb}/citypes", GetCITypes).Methods("GET")
	router.HandleFunc("/cmdbs/{cmdb}/citypes", AddCIType).Methods("POST")
	router.HandleFunc("/cmdbs/{cmdb}/citypes/{name}", GetCITypeByName).Methods("GET")
	router.HandleFunc("/cmdbs/{cmdb}/citypes/{name}", DeleteCITypeByName).Methods("DELETE")

	// CI routes
	router.HandleFunc("/cmdbs/{cmdb}/{citype}", GetCIs).Methods("GET")
	router.HandleFunc("/cmdbs/{cmdb}/{citype}", AddCI).Methods("POST")
	router.HandleFunc("/cmdbs/{cmdb}/{citype}/{id}", GetCIById).Methods("GET")
	router.HandleFunc("/cmdbs/{cmdb}/{citype}/{id}", DeleteCIById).Methods("DELETE")

	// Init Negroni
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), NewAuthHandler())
	n.UseHandler(router)

	return n
}

func Serve() {
	// Get configuration
	config, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	n := GetServer()
	n.Run(fmt.Sprintf("%s:%d", config.Server.ListenOn, config.Server.ListenPort))
}

func ExpandPath(path string) string {
	if path[:1] == "~" {
		usr, _ := user.Current()
		path = fmt.Sprintf("%s%s", usr.HomeDir, path[1:])
	}

	return path
}

func IsValidShortName(name string) bool {
	match, err := regexp.MatchString("^[a-zA-Z0-9-_]+$", name)
	if err != nil {
		log.Panic(err)
	}
	return match
}
