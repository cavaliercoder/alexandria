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
 * along with this progradb.  If not, see <http://www.gnu.org/licenses/>.
 * package controllers
 */
package controllers

import (
    "alexandria/api/services"
    "alexandria/api/models"
    
    "log"
    "net/http"
    
    "github.com/go-martini/martini"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

const (
    database string = "alexandria"
    collection string = "config"    
)

type ConfigController struct {
    BaseController
}

func (c ConfigController) Init(r martini.Router)  error {
    // Add routes
    r.Get("/config", c.getConfig)
    r.Post("/config/actions/initialize", c.initConfig)
    r.Post("/config/actions/destroy", c.clearConfig)
    return nil
}

func (c ConfigController) getConfig(dbsession *services.Database, r *services.Renderer) {
    var config models.Config
    err := dbsession.DB(database).C(collection).Find(nil).One(&config)
    r.Handle(err)
    
    r.Render(http.StatusOK, config)
}

func (c ConfigController) initConfig(dbsession *services.Database, r *services.Renderer) {
    log.Printf("Received request to initialize API configuration from %s", r.RemoteAddr)
        
    // Does configuration exist?
    db := dbsession.DB(database)
    count, err := db.C("config").Find(bson.M{}).Count()
    if err != nil { log.Panic(err) }
    
    if count != 0 {
        log.Printf("Configuration already initialized")
        r.NotFound()  
    } else {
        // Create default configuration
        log.Print("Initializing API configuration")
        
        // Create root tenant
        db.C("tenants").Create(&mgo.CollectionInfo{})
        db.C("tenants").EnsureIndex(mgo.Index{ Key: []string{"code"}, Unique: true})
        
        tenant := models.Tenant{}
        tenant.Init()
        tenant.Name = "Default Tenant"
        
        err = db.C("tenants").Insert(tenant)
        if err != nil { log.Fatal(err) }
        log.Printf("Created root tenant with ID %s", tenant.Id.Hex())
        
        // Create root user   
        db.C("users").Create(&mgo.CollectionInfo{})
        db.C("users").EnsureIndex(mgo.Index{ Key: []string{"email"}, Unique: true})
        db.C("users").EnsureIndex(mgo.Index{ Key: []string{"apiKey"}, Unique: true, Sparse: true})
    
        user := models.User{}
        user.Init()
        user.Email = "root"
        user.TenantId = tenant.Id
        
        err = db.C("users").Insert(user)
        if err != nil { log.Fatal(err) }
        log.Printf("Created root user with ID %s", user.Id.Hex())
        
        // Create configuration
        db.C("config").Create(&mgo.CollectionInfo{})
    
        config := models.Config{}
        config.Init()
        config.Initialized = true
        config.RootTenant = tenant.Id
        config.RootUser = user.Id
        config.Version = "1.0.0"
            
        err = db.C("config").Insert(config)
        if err != nil { log.Fatal(err) }
        log.Printf("Configuration initialization completed successfully")
        
        r.ResponseWriter.Header().Set("Location", "/config")
        r.Render(http.StatusCreated, config)
    }
}

func (c ConfigController) clearConfig(dbsession *services.Database, r *services.Renderer) {
    err := dbsession.DB(database).DropDatabase()
    if err != nil {
        r.Handle(err)
    }
    
    r.Render(200, "Death from above!!!")
}