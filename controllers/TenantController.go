/*
 * Alexandria CMDB - Open source common.management database
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
	"github.com/cavaliercoder/alexandria/models"
	"github.com/cavaliercoder/alexandria/services"
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

func (c *TenantController) GetPath() string {
    return "/tenants"
}

func (c *TenantController) InitRoutes(r martini.Router) error {
	r.Get("/", c.getTenants)
	r.Post("/", binding.Bind(models.Tenant{}), c.addTenant)
	r.Get("/:id", c.getTenant)

	return nil
}

func (c *TenantController) getTenant(r *services.ApiContext, params martini.Params) {
	var tenant models.Tenant
	id := strings.ToUpper(params["id"])
	err := r.DB.GetOneById("tenants", id, &tenant)

	r.Handle(err)

	r.Render(http.StatusOK, tenant)
}

func (c *TenantController) getTenants(r *services.ApiContext) {
	var tenants []models.Tenant
	err := r.DB.GetAll("tenants", nil, &tenants)
	r.Handle(err)

	r.Render(http.StatusOK, tenants)
}

func (c *TenantController) addTenant(tenant models.Tenant, r *services.ApiContext) {
	tenant.Init(r.DB.NewId())
	err := r.DB.Insert("tenants", &tenant)
	if err != nil {
		log.Fatal(err)
	}

	r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/tenants/%s", r.DB.IdToString(tenant.Id)))
	r.Render(http.StatusCreated, "")
}
