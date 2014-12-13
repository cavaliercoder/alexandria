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
    "alexandria/api/application"
    "alexandria/api/models"
    
    "log"
    "net/http"
    
    "github.com/go-martini/martini"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type ConfigController struct {
    BaseController
}

func (c ConfigController) Init(app *application.AppContext)  error {
    c.app = app
        
    // Add routes
    c.app.Martini.Get("/config", c.GetConfig)
    
    // Initialize database
    c.app.Db.C("config").Create(&mgo.CollectionInfo{})
    
    return nil
}

func (c ConfigController) GetConfig(params martini.Params, r *http.Request, w http.ResponseWriter) {
    if r.URL.Query().Get("init") == "true" {
        log.Printf("Received request to initialize API configuration from %s", r.RemoteAddr)
    
        if ! c.initDb() {
            log.Printf("Configuration already initialized. Sending 404") // Or "invalid parameter???"
            w.WriteHeader(http.StatusNotFound)
            return
        }
    }
    
    c.GetEntities("config", models.Config{}, nil, w)
}

func (c ConfigController) initDb() bool {
    // Does configuration exist?
    count, err := c.app.Db.C("config").Find(bson.M{}).Count()
    if err != nil { log.Panic(err) }    
    if count == 0 {
        // Create default configuration
        log.Print("Initializing API configuration")
        
        // Create root tenant
        tenant := models.Tenant{
            Name: "Root tenant",
        }        
        tenant.SetCreated()
        
        err = c.app.Db.C("tenants").Insert(tenant)
        if err != nil { log.Fatal(err) }
        log.Printf("Created root tenant with ID %s", tenant.Id)
        
        // Create root user
        user := models.User{
            Email: "root",
            TenantId: tenant.Id,
        }
        user.SetCreated()
        
        err = c.app.Db.C("users").Insert(user)
        if err != nil { log.Fatal(err) }
        log.Printf("Created root user with ID %s", user.Id)
        
        // Create configuration
        config := models.Config{
            Initialized: true,
            RootTenant: tenant.Id,
            RootUser: user.Id,
        }
        config.SetCreated()
        
        err = c.app.Db.C("config").Insert(config)
        if err != nil { log.Fatal(err) }
        log.Printf("Configuration initialization completed successfully")
        
        return true
    }
    
    return false
}