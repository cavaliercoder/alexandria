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
    "crypto/sha1"
    "log"
    "encoding/json"
    "fmt"
    "net/http"
    
    "alexandria/api/application"
    "alexandria/api/models"
    
    "github.com/go-martini/martini"
    "github.com/martini-contrib/binding"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type UserController struct {
    BaseController
}

func (c UserController) Init(app *application.AppContext)  error {
    c.app = app
    
    // Add routes
    c.app.Martini.Get("/users", c.GetUsers)
    c.app.Martini.Get("/users/:email", c.GetUserByEmail)
    c.app.Martini.Post("/users", binding.Bind(models.User{}), c.AddUser)
  
    // Initialize database
    c.app.Db.C("users").Create(&mgo.CollectionInfo{})
    c.app.Db.C("users").EnsureIndex(mgo.Index{ Key: []string{"Email"}, Unique: true})
    c.app.Db.C("users").EnsureIndex(mgo.Index{ Key: []string{"apiKey"}, Unique: true, Sparse: true})
    
    return nil
}

func (c UserController) GetUsers(w http.ResponseWriter) {    
    c.GetEntities("users", models.User{}, nil, w)
}

func (c UserController) GetUserByEmail(params martini.Params, w http.ResponseWriter) {
    var user models.User
    err := c.app.Db.C("users").Find(bson.M{"email": params["email"]}).One(&user)
    c.Handle(err)
    
    c.RenderJson(w, user)
}

func (c UserController) AddUser(user models.User, w http.ResponseWriter) {
    // Make sure user doesn't already exist
    count, err := c.app.Db.C("users").Find(bson.M{"email": user.Email}).Count()
    c.Handle(err)    
    if count > 0 {
        w.WriteHeader(http.StatusConflict)
        log.Panic(fmt.Sprintf("A user account already exists for email %s", user.Email))
    }
    
    // Create API key
    jsonHash, err := json.Marshal(user)
    c.Handle(err)
    shaSum := sha1.Sum(jsonHash)
    user.ApiKey = fmt.Sprintf("%x", shaSum)
    
    c.AddEntity("users", fmt.Sprintf("/users/%s", user.Email), &user, w)
}