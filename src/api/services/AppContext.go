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

/*
import (
        "encoding/json"
        "log"
        "net/http"
        "reflect"
        
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
        
        "alexandria/api/models"
)

type AppContext struct {
    Request         *http.Request
    ResponseWriter  http.ResponseWriter
    MongoSession    *mgo.Session
    Params          *martini.Params
}

// Wire the context to a martini.Handler service
func AppContextService() martini.Handler {
    
    return func(req *http.Request, res http.ResponseWriter, c martini.Context, params *martini.Params) {
        appContext := AppContext{
            Request:        req,
            ResponseWriter: res,
            MongoSession:   mgoSession.Clone(),
            Params:         params,
        }
        
        defer appContext.MongoSession.Close()
        
        c.Next()
    }
}

func (c AppContext) Handle(err error) {
    if err != nil {
        log.Panic(err)
    }
}

func (c AppContext) GetEntities(collection string, resultType interface{}, query interface{}) {
    //if query == nil { query = bson.M{} }
    
    typ := reflect.TypeOf(resultType)
    results := reflect.New(reflect.SliceOf(typ)).Interface()
    
    dbcollection := c.MongoSession.DB("alexandria").C(collection)
    err := dbcollection.Find(query).All(results)
    c.Handle(err)
    
    c.RenderJson(results)
}

func (c AppContext) AddEntity(collection string, uri string, model interface{}) {
    baseModel, success := model.(models.ModelInterface)
    if ! success {
        log.Panic("Model is invalid")
        return
    }
    
    // Insert new entity
    baseModel.SetCreated()    
    err := c.MongoSession.DB("alexandria").C(collection).Insert(baseModel)
    c.Handle(err)
    
    // Update response headers
    c.ResponseWriter.Header().Set("Location", uri)
    c.ResponseWriter.WriteHeader(http.StatusCreated)
    
    c.RenderJson(model)
}
*/