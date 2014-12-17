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
    "log"
    "fmt"
    "net/http"
    
    "alexandria/api/services"
    "alexandria/api/models"
    
    "github.com/go-martini/martini"
    "github.com/martini-contrib/binding"
    "gopkg.in/mgo.v2/bson"
)

type UserController struct {
    BaseController
}

func (c UserController) Init(r martini.Router)  error {
    
    // Add routes
    r.Get("/users", c.getUsers)
    r.Get("/users/:email", c.getUserByEmail)
    r.Post("/users", binding.Bind(models.User{}), c.addUser)
    
    return nil
}

func (c UserController) getUsers(dbsession *services.Database, r *services.Renderer) {    
    var users []models.User
    err := dbsession.DB("alexandria").C("users").Find(nil).All(&users)
    r.Handle(err)
    
    r.Render(http.StatusOK, users)
}

func (c UserController) getUserByEmail(dbsession *services.Database, r *services.Renderer, params martini.Params) {
    var user models.User
    err := dbsession.DB("alexandria").C("users").Find(bson.M{"email": params["email"]}).One(&user)
    r.Handle(err)
    
    r.Render(http.StatusOK, user)
}

func (c UserController) addUser(user models.User, dbsession *services.Database, r *services.Renderer) {
    user.Init()
    err := dbsession.DB("alexandria").C("users").Insert(user)
    if err != nil { log.Fatal(err) }
    
    r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/users/%s", user.Email))
    r.Render(http.StatusCreated, "")
}