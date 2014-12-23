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
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"

	"github.com/cavaliercoder/alexandria/database"
	"github.com/cavaliercoder/alexandria/models"
)

type ApiContext struct {
	*http.Request
	http.ResponseWriter
	context  martini.Context // Martini context
	DB       database.Driver // Database driver
	AuthUser *models.User    // Authenticated user
	AuthTenant *models.Tenant // Authenticated user's tenancy
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

		// Get authenticated user
		user, tenant, err := getAuth(req, clone)
		if err != nil {
			log.Printf(err.Error())
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		
		// Create context
		r := &ApiContext{req, res, c, clone, user, tenant}

		// Wire it up
		c.Map(r)
		c.Next()
	}
}

func getAuth(req *http.Request, db database.Driver) (*models.User, *models.Tenant, error) {
	var user models.User
	var tenant models.Tenant

	// Get API key from request header
	apiKey := req.Header.Get("X-Auth-Token")
	if apiKey == "" {
		return nil, nil, errors.New("Authentication token not set")
	}

	// Get user account
	err := db.GetOne("users", database.M{"apikey": apiKey}, &user)
	if err != nil {
		return nil, nil, err
	}
	
	// Get user tenancy
	err = db.GetOneById("tenants", user.TenantId, &tenant)
	if err != nil {
		return nil, nil, err
	}

	return &user, &tenant, nil
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
