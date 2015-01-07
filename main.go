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
 */
package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
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

	// Initr public routes
	pub := mux.NewRouter().PathPrefix(ApiV1Prefix).Subrouter()
	pub.HandleFunc("/info", GetApiInfo).Methods("GET")
	pub.HandleFunc("/apikey", GetApiKey).Methods("POST")

	// Init private routes
	priv := mux.NewRouter().PathPrefix(ApiV1Prefix).Subrouter()

	// User routes
	priv.HandleFunc("/users", GetUsers).Methods("GET")
	priv.HandleFunc("/users", AddUser).Methods("POST")
	priv.HandleFunc("/users/current", GetCurrentUser).Methods("GET")
	priv.HandleFunc("/users/{email}", GetUserByEmail).Methods("GET")
	priv.HandleFunc("/users/{email}", DeleteUserByEmail).Methods("DELETE")
	priv.HandleFunc("/users/{email}/password", SetUserPassword).Methods("PATCH")

	// Tenant routes
	priv.HandleFunc("/tenants", GetTenants).Methods("GET")
	priv.HandleFunc("/tenants", AddTenant).Methods("POST")
	priv.HandleFunc("/tenants/current", GetCurrentTenant).Methods("GET")
	priv.HandleFunc("/tenants/{code}", GetTenantByCode).Methods("GET")
	priv.HandleFunc("/tenants/{code}", DeleteTenantByCode).Methods("DELETE")

	// CMDB routes
	priv.HandleFunc("/cmdbs", GetCmdbs).Methods("GET")
	priv.HandleFunc("/cmdbs", AddCmdb).Methods("POST")
	priv.HandleFunc("/cmdbs/{name}", GetCmdbByName).Methods("GET")
	priv.HandleFunc("/cmdbs/{name}", DeleteCmdbByName).Methods("DELETE")

	// CI Type routes
	priv.HandleFunc("/cmdbs/{cmdb}/citypes", GetCITypes).Methods("GET")
	priv.HandleFunc("/cmdbs/{cmdb}/citypes", AddCIType).Methods("POST")
	priv.HandleFunc("/cmdbs/{cmdb}/citypes/{name}", GetCITypeByName).Methods("GET")
	priv.HandleFunc("/cmdbs/{cmdb}/citypes/{name}", DeleteCITypeByName).Methods("DELETE")

	// CI routes
	priv.HandleFunc("/cmdbs/{cmdb}/{citype}", GetCIs).Methods("GET")
	priv.HandleFunc("/cmdbs/{cmdb}/{citype}", AddCI).Methods("POST")
	priv.HandleFunc("/cmdbs/{cmdb}/{citype}/{id}", GetCIById).Methods("GET")
	priv.HandleFunc("/cmdbs/{cmdb}/{citype}/{id}", DeleteCIById).Methods("DELETE")

	// Init Negroni with public routes
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(pub)

	// If the public router can't find a match, pass the request to the
	// private Negroni instance
	npriv := negroni.New(NewAuthHandler())
	npriv.UseHandler(priv)
	pub.NotFoundHandler = npriv

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
