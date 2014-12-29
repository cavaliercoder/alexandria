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
		c.fail(res, req)
		return
	} else {
		context := GetAuthContext(req)
		if context == nil {
			log.Printf("No user or tenancy found with API Key: %s", apiKey)
			c.fail(res, req)
			return
		}
	}

	// Process request chain
	next(res, req)

	// Remove the user from the request cache
	delete(authCache, req)
}

func (c *AuthHandler) fail(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusUnauthorized)
	res.Write([]byte("401 Unauthorized"))
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
