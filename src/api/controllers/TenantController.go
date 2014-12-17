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
	"alexandria/api/models"
	"alexandria/api/services"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strings"
)

type TenantController struct {
	BaseController
}

func (c TenantController) Init(r martini.Router) error {
	// Add routes
	r.Get("/tenants", c.getTenants)
	r.Get("/tenants/:id", c.getTenant)
	r.Post("/tenants", binding.Bind(models.Tenant{}), c.addTenant)

	return nil
}

func (c TenantController) getTenant(dbsession *services.Database, r *services.Renderer, params martini.Params) {
	var tenant models.Tenant
	code := strings.ToUpper(params["id"])
	err := dbsession.DB("alexandria").C("tenants").Find(bson.M{"code": code}).One(&tenant)

	r.Handle(err)

	r.Render(http.StatusOK, tenant)
}

func (c TenantController) getTenants(dbsession *services.Database, r *services.Renderer) {
	var tenants []models.Tenant
	err := dbsession.DB("alexandria").C("tenants").Find(nil).All(&tenants)
	r.Handle(err)

	r.Render(http.StatusOK, tenants)
}

func (c TenantController) addTenant(tenant models.Tenant, dbsession *services.Database, r *services.Renderer) {
	tenant.Init()
	err := dbsession.DB("alexandria").C("tenants").Insert(tenant)
	if err != nil {
		log.Fatal(err)
	}

	r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/tenants/%s", tenant.Code))
	r.Render(http.StatusCreated, "")
}
