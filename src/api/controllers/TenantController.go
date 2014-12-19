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
	"alexandria/api/database"
	"alexandria/api/models"
	"alexandria/api/services"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"log"
	"net/http"
	"strings"
)

type TenantController struct {
	controller
}

func (c *TenantController) Init(r martini.Router) error {
	// Add routes
	r.Get("/tenants", c.getTenants)
	r.Get("/tenants/:id", c.getTenant)
	r.Post("/tenants", binding.Bind(models.Tenant{}), c.addTenant)

	return nil
}

func (c *TenantController) getTenant(dbdriver database.Driver, r *services.ApiContext, params martini.Params) {
	var tenant models.Tenant
	code := strings.ToUpper(params["id"])
	err := dbdriver.GetOne("tenants", database.M{"code": code}, &tenant)

	r.Handle(err)

	r.Render(http.StatusOK, tenant)
}

func (c *TenantController) getTenants(dbdriver database.Driver, r *services.ApiContext) {
	var tenants []models.Tenant
	err := dbdriver.GetAll("tenants", nil, &tenants)
	r.Handle(err)

	r.Render(http.StatusOK, tenants)
}

func (c *TenantController) addTenant(tenant models.Tenant, dbdriver database.Driver, r *services.ApiContext) {
	tenant.Init()
	err := dbdriver.Insert("tenants", tenant)
	if err != nil {
		log.Fatal(err)
	}

	r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/tenants/%s", tenant.Code))
	r.Render(http.StatusCreated, "")
}
