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
    "github.com/go-martini/martini"
    "github.com/martini-contrib/binding"
    "gopkg.in/mgo.v2"
    "net/http"
    "alexandria/api/models"
    "fmt"
)

type tenantController struct {
    baseController
}

func NewTenantController(m *martini.ClassicMartini, db *mgo.Database) (*tenantController, error) {
    c := new(tenantController)    
    c.m = m
    c.db = db    
    
    // Add routes
    m.Get("/tenants", c.GetTenants)
    //m.Get("/tenants/:email", c.GetUserByEmail)
    m.Post("/tenants", binding.Bind(models.Tenant{}), c.AddTenant)
    
    return c, nil
}

func (c tenantController) GetTenants(w http.ResponseWriter) {    
    c.GetEntities("tenants", w)
}

func (c tenantController) AddTenant(tenant models.Tenant, w http.ResponseWriter) {
    c.AddEntity("tenants", fmt.Sprintf("/tenants/%s", tenant.Id), &tenant, w)
}