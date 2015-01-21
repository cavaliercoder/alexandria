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
 */
package main

import (
	"errors"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type AuthHandler struct {
}

type AuthContext struct {
	User   *User
	Tenant *Tenant
}

type AuthMap map[*http.Request]*AuthContext

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (c *AuthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Get API key from request header
	apiKey := req.Header.Get("X-Auth-Token")
	if apiKey == "" {
		log.Printf("X-Auth-Token header not set")
		ErrUnauthorized(res, req)
		return
	} else {
		context := GetAuthContext(req)
		if context == nil {
			log.Printf("No user or tenancy found with API Key: %s", apiKey)
			ErrUnauthorized(res, req)
			return
		}
	}

	// Process request chain
	next(res, req)

	// Remove the user from the request cache
	delete(authCache, req)
}

var authCache AuthMap

func GetAuthContext(req *http.Request) *AuthContext {
	// Initialize the context cache
	if authCache == nil {
		authCache = AuthMap{}
	}

	// Is the user cached already?
	context, ok := authCache[req]
	if ok {
		return context
	}

	// Get API key from request header
	apiKey := req.Header.Get("X-Auth-Token")
	if apiKey == "" {
		return nil
	} else {
		// Find the user
		var user User
		err := RootDb().C("users").Find(M{"apikey": apiKey}).One(&user)
		if err == mgo.ErrNotFound {
			return nil
		} else if err != nil {
			log.Printf("Error retrieving API user from the database: %s", err.Error())
			return nil
		}

		// Find the tenant
		var tenant Tenant
		err = RootDb().C("tenants").FindId(user.TenantId).One(&tenant)
		if err == mgo.ErrNotFound {
			return nil
		} else if err != nil {
			log.Printf("Error retrieving API tenant from the database: %s", err.Error())
			return nil
		}

		// Add the conext to the cache
		context = &AuthContext{&user, &tenant}
		authCache[req] = context
		return context
	}
}

// GetApiKey accepts a JSON request body with a user name and password
// encapsulated and returns the user's API key
func GetApiKey(res http.ResponseWriter, req *http.Request) {
	// Parse the request body. Should be:
	// {
	//    "username":"some@email.com",
	//    "password":"S0m3P4ssw0RD"
	// }
	body := make(map[string]string)
	err := Bind(req, &body)
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}
	if body["username"] == "" || body["password"] == "" {
		err = errors.New("Username or password not specified")
		ErrBadRequest(res, req, err)
		return
	}

	// Find the user account
	var user User
	err = RootDb().C("users").Find(M{"email": body["username"]}).One(&user)
	if err != nil {
		ErrUnauthorized(res, req)
		return
	}

	// Validate the password
	if !CheckPassword(user.PasswordHash, body["password"]) {
		ErrUnauthorized(res, req)
		return
	}

	// Formulate response
	key := map[string]string{
		"apiKey": user.ApiKey,
	}

	Render(res, req, http.StatusOK, &key)
}
