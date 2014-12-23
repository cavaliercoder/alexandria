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
package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cavaliercoder/alexandria/database"
	"github.com/cavaliercoder/alexandria/models"
	"github.com/cavaliercoder/alexandria/services"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

type UserController struct {
	controller
}

func (c *UserController) GetPath() string {
    return "/users"
}

func (c *UserController) InitRoutes(r martini.Router) error {
	r.Get("/", c.getUsers)
	r.Post("/", binding.Bind(models.User{}), c.addUser)
	r.Get("/:email", c.getUserByEmail)
	r.Delete("/:email", c.deleteUserByEmail)

	return nil
}

func (c *UserController) getUsers(r *services.ApiContext) {
	var users []models.User
	err := r.DB.GetAll("users", nil, &users)
	r.Handle(err)

	r.Render(http.StatusOK, users)
}

func (c *UserController) getUserByEmail(r *services.ApiContext, params martini.Params) {
        // TODO: route the correct database to the controllers
        // TODO: allow for adding/removing databases from tenants
        
	var user models.User
	err := r.DB.GetOne("users", database.M{"email": params["email"]}, &user)
	if r.Handle(err) {
		return
	}

	r.Render(http.StatusOK, user)
}

func (c *UserController) addUser(user models.User, r *services.ApiContext) {
	user.Init(r.DB.NewId())
	user.TenantId = r.AuthUser.TenantId

	err := r.DB.Insert("users", &user)
	if err != nil {
		log.Panic(err)
	}

	r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/users/%s", user.Email))
	r.Render(http.StatusCreated, "")
}

func (c *UserController) deleteUserByEmail(r *services.ApiContext, params martini.Params) {
	err := r.DB.Delete("users", database.M{"tenantid": r.AuthUser.TenantId, "email": params["email"]})
	if r.Handle(err) {
		return
	}

	r.Render(http.StatusNoContent, "")
}
