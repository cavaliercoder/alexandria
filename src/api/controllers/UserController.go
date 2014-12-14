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
/*
import (
    "crypto/sha1"
    "log"
    "encoding/json"
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
    r.Get("/users", c.GetUsers)
    r.Get("/users/:email", c.GetUserByEmail)
    r.Post("/users", binding.Bind(models.User{}), c.AddUser)
    
    return nil
}

func (c UserController) GetUsers(context *services.AppContext) {    
    context.GetEntities("users", models.User{}, nil)
}

func (c UserController) GetUserByEmail(context *services.AppContext) {
    var user models.User
    err := context.MongoSession.DB("alexandria").C("users").Find(bson.M{"email": (*context.Params)["email"]}).One(&user)
    context.Handle(err)
    
    context.RenderJson(user)
}

func (c UserController) AddUser(user models.User, context *services.AppContext) {
    // Make sure user doesn't already exist
    count, err := context.MongoSession.DB("alexandria").C("users").Find(bson.M{"email": user.Email}).Count()
    context.Handle(err)    
    if count > 0 {
        context.ResponseWriter.WriteHeader(http.StatusConflict)
        log.Panic(fmt.Sprintf("A user account already exists for email %s", user.Email))
    }
    
    // Create API key
    jsonHash, err := json.Marshal(user)
    context.Handle(err)
    shaSum := sha1.Sum(jsonHash)
    user.ApiKey = fmt.Sprintf("%x", shaSum)
    
    context.AddEntity("users", fmt.Sprintf("/users/%s", user.Email), &user)
}
*/