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
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
		_, err = GetConfigFromFile(confFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	// Init Mux routes
	router := mux.NewRouter()
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users", AddUser).Methods("POST")
	router.HandleFunc("/users/{email}", GetUserByEmail).Methods("GET")

		
	// Init Negroni      
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
	
	// Git'er done
	log.Printf("Initialization complete")
}

func Handle(err error) bool {
	if err != nil {
		log.Panic(err)
		return true
	}
	
	return false
}

func Render(res http.ResponseWriter, req *http.Request, status int, v interface{}) {
	format := req.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	switch format {
	case "json":
		RenderJson(res, req, status, v)

	default:
		log.Panic(fmt.Sprintf("Unsupported output format: %s", format))
	}
}

func RenderJson(res http.ResponseWriter, req *http.Request, status int, v interface{}) {
	if v == nil {
		v = new(struct{})
	}

	var err error
	var data []byte
	if req.URL.Query().Get("pretty") == "true" {
		data, err = json.MarshalIndent(v, "", "    ")
	} else {
		data, err = json.Marshal(v)
	}
	if err != nil {
		log.Panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	res.Write(data)
}

func Bind(req *http.Request, v interface{}) error {
	if req.Body != nil {
		defer req.Body.Close()
		
		err := json.NewDecoder(req.Body).Decode(v)
		
		if err != nil && err != io.EOF {
			return err
		}
	}
	
	return nil
}

func GetPathVar(req *http.Request, name string) string {
	vars := mux.Vars(req)
	result := vars[name]
	
	if name == "" {
		log.Panic(fmt.Sprintf("No such variable declared: %s", name))
	}
	
	return result
}