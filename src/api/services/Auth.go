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
	"net/http"

	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"

	"alexandria/api/models"
)

// validate an api key
func ApiKeyValidation() martini.Handler {
	db := DbConnect()

	return func(res http.ResponseWriter, req *http.Request) {
		apiKey := req.Header.Get("X-Auth-Token")
		if apiKey == "" {
			res.WriteHeader(http.StatusUnauthorized)
		} else {
			var user models.User
			var tenant models.Tenant
			err := db.DB("alexandria").C("users").Find(bson.M{"apiKey": apiKey}).One(&user)
			if err != nil {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}

			err = db.DB("alexandria").C("tenants").FindId(user.TenantId).One(&tenant)
			if err != nil {
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}
}
