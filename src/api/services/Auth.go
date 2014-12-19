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
	"log"
	"net/http"

	"github.com/go-martini/martini"

	"alexandria/api/database"
	"alexandria/api/models"
)

// validate an api key
func ApiKeyValidation() martini.Handler {
	db, err  := database.Connect()
	if err != nil { log.Panic(err) }
	
	return func(res http.ResponseWriter, req *http.Request) {
		apiKey := req.Header.Get("X-Auth-Token")
		if apiKey == "" {
			res.WriteHeader(http.StatusUnauthorized)
		} else {
			var user models.User
			var tenant models.Tenant
			err := db.GetOne("users", database.M{"apikey": apiKey}, &user)
			if err != nil {
				log.Print(err)
				res.WriteHeader(http.StatusUnauthorized)
				return
			}

			err = db.GetOneById("tenants", user.TenantId, &tenant)
			if err != nil {
				log.Print(err)
				res.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}
}
