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
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (c *AuthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Get API key from request header
	apiKey := req.Header.Get("X-Auth-Token")
	if apiKey == "" {
		log.Printf("X-Auth-Token header not set")
		c.fail(res, req)
	} else {

		// Find the user
		var user User
		err := RootDb().C("users").Find(M{"apikey": apiKey}).One(&user)
		if err == mgo.ErrNotFound {
			log.Printf("No user found with API Key %s", apiKey)
			c.fail(res, req)
		}
	}

	next(res, req)
}

func (c *AuthHandler) fail(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusUnauthorized)
	res.Write([]byte("401 Unauthorized"))
}
