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
    "gopkg.in/mgo.v2/bson"
    "encoding/json"
    "net/http"
    "log"
    "alexandria/api/application"
    "alexandria/api/models"
)

type ControllerInterface interface {
    Init(app *application.AppContext)      error
}

type BaseController struct {
    app *application.AppContext
}

func (c BaseController) RenderJson(w http.ResponseWriter, v interface{}) {
    if v == nil {
        var empty []struct{}
        v = empty
    }
    
    json, err := json.Marshal(v)
    c.Handle(err)
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(json)
}

func (c BaseController) Handle(err error) {
    if err != nil {
        log.Panic(err)
    }
}

func (c BaseController) GetEntities(collection string, w http.ResponseWriter) {    
    dbcollection := c.app.Db.C(collection)
    var results []interface{}
    err := dbcollection.Find(bson.M{}).All(&results)
    c.Handle(err)
    
    c.RenderJson(w, results)
}

func (c BaseController) AddEntity(collection string, uri string, model interface{}, w http.ResponseWriter) {
    baseModel, success := model.(models.ModelInterface)
    if ! success {
        log.Panic("Model is invalid")
        return
    }
    
    // Insert new user
    baseModel.SetCreated()    
    err := c.app.Db.C(collection).Insert(baseModel)
    c.Handle(err)
    
    // Update headers
    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Location", uri)
    
    c.RenderJson(w, model)
}