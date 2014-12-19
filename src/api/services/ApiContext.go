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
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"

	"alexandria/api/database"
	"alexandria/api/models"
)

type ApiContext struct {
	*http.Request
	http.ResponseWriter
	context  martini.Context // Martini context
	DB       database.Driver // Database driver
	AuthUser *models.User    // Authenticated user
}

// Wire the service
func ApiContextService() martini.Handler {
	db, err := database.Connect()
	if err != nil {
		log.Panic(err)
	}

	return func(req *http.Request, res http.ResponseWriter, c martini.Context) {
		// Connect to the database
		clone, err := db.Clone()
		if err != nil {
			log.Panic(err)
		}
		defer clone.Close()

		// Create context
		r := &ApiContext{req, res, c, clone, nil}

		// Get authenticated user
		user, err := r.GetAuthUser()
		if r.Handle(err) {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		r.AuthUser = user

		// Wire it up
		c.Map(r)
		c.Next()
	}
}

func (c *ApiContext) GetAuthUser() (*models.User, error) {
	var user models.User

	apiKey := c.Request.Header.Get("X-Auth-Token")
	if apiKey == "" {
		return nil, errors.New("Authentication token not set")
	}

	err := c.DB.GetOne("users", database.M{"apikey": apiKey}, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *ApiContext) Handle(err error) bool {
	if err == nil {
		return false
	}

	switch err.Error() {
	case "not found":
		c.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.ResponseWriter.Write([]byte("404 page not found\n"))
	default:
		log.Panic(err)
	}

	return true
}

func (c *ApiContext) NotFound() {
	c.ResponseWriter.WriteHeader(http.StatusNotFound)
}

func (c *ApiContext) Render(status int, v interface{}) {
	format := c.Request.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	switch format {
	case "json":
		c.JSON(status, v)

	default:
		log.Panic(fmt.Sprintf("Unsupported output format: %s", format))
	}
}

func (c *ApiContext) JSON(status int, v interface{}) {
	if v == nil {
		v = new(struct{})
	}

	var err error
	var data []byte
	if c.Request.URL.Query().Get("pretty") == "true" {
		data, err = json.MarshalIndent(v, "", "    ")
	} else {
		data, err = json.Marshal(v)
	}
	if err != nil {
		log.Panic(err)
	}

	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(status)

	c.ResponseWriter.Write(data)
}
