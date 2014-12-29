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
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	model     `json:"-" bson:",inline"`
	TenantId  interface{} `json:"-"`
	ApiKey    string      `json:"-"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     string      `json:"email"`
	Password  string      `json:"-"`
}

func (c *User) InitModel() {
	c.model.InitModel()
	c.GenerateApiKey()
}

func (c *User) GenerateApiKey() string {
	// Create API key
	jsonHash, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	shaSum := sha1.Sum(jsonHash)
	c.ApiKey = fmt.Sprintf("%x", shaSum)

	return c.ApiKey
}

func GetUsers(res http.ResponseWriter, req *http.Request) {
	var users []User
	err := RootDb().C("users").Find(nil).All(&users)
	Handle(res, req, err)

	Render(res, req, http.StatusOK, users)
}

func GetUserByEmail(res http.ResponseWriter, req *http.Request) {
	email := GetPathVar(req, "email")

	var user User
	err := RootDb().C("users").Find(M{"email": email}).One(&user)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusOK, user)
}

func AddUser(res http.ResponseWriter, req *http.Request) {
	var user User
	err := Bind(req, &user)
	if Handle(res, req, err) {
		return
	}

	user.InitModel()
	//user.TenantId = r.AuthUser.TenantId

	err = RootDb().C("users").Insert(&user)
	if Handle(res, req, err) {
		return
	}

	RenderCreated(res, req, fmt.Sprintf("%s/users/%s", ApiV1Prefix, user.Email))
}

func DeleteUserByEmail(res http.ResponseWriter, req *http.Request) {
	email := GetPathVar(req, "email")

	// TODO: Ensure only users for current tenant can be deleted
	err := RootDb().C("users").Remove(M{"email": email})
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusNoContent, "")
}
