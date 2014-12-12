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
    "github.com/martini-contrib/binding"
    "net/http"
    "alexandria/api/application"
    "alexandria/api/models"
    "fmt"
)

type TenantController struct {
    BaseController
}

func (c TenantController) Init(app *application.AppContext)  error {
    c.app = app
        
    // Add routes
    c.app.Martini.Get("/tenants", c.GetTenants)
    c.app.Martini.Post("/tenants", binding.Bind(models.Tenant{}), c.AddTenant)
    
    return nil
}

func (c TenantController) GetTenants(w http.ResponseWriter) {    
    c.GetEntities("tenants", w)
}

func (c TenantController) AddTenant(tenant models.Tenant, w http.ResponseWriter) {
    c.AddEntity("tenants", fmt.Sprintf("/tenants/%s", tenant.Id), &tenant, w)
}