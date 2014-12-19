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
package controllers

import (
	"fmt"
	"log"
	"net/http"

	"alexandria/api/database"
	"alexandria/api/models"
	"alexandria/api/services"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

type UserController struct {
	controller
}

func (c *UserController) Init(r martini.Router) error {

	// Add routes
	r.Get("/users", c.getUsers)
	r.Get("/users/:email", c.getUserByEmail)
	r.Post("/users", binding.Bind(models.User{}), c.addUser)

	return nil
}

func (c *UserController) getUsers(dbdriver database.Driver, r *services.Renderer) {
	var users []models.User
	err := dbdriver.GetAll("users", nil, &users)
	r.Handle(err)

	r.Render(http.StatusOK, users)
}

func (c *UserController) getUserByEmail(dbdriver database.Driver, r *services.Renderer, params martini.Params) {
	var user models.User
	err := dbdriver.GetOne("users", database.M{"email": params["email"]}, &user)
	if r.Handle(err) { return }

	r.Render(http.StatusOK, user)
}

func (c *UserController) addUser(user models.User, dbdriver database.Driver, r *services.Renderer) {
	user.Init()
        err := dbdriver.Insert("users", user)
	if err != nil { log.Panic(err) }

	r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/users/%s", user.Email))
	r.Render(http.StatusCreated, "")
}
