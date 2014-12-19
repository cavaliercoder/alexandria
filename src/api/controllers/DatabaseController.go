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
	"fmt"
	"log"
	"net/http"

	"alexandria/api/database"
	"alexandria/api/models"
	"alexandria/api/services"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

type DatabaseController struct {
	controller
}

func (c *DatabaseController) Init(r martini.Router) error {
	// Add routes
	r.Get("/databases", c.getDatabases)
	r.Get("/databases/:shortname", c.getDatabaseByShortName)
	r.Post("/databases", binding.Bind(models.Database{}), c.createDatabase)

	return nil
}

func (c *DatabaseController) getDatabases(r *services.ApiContext) {
	var databases []models.Database
	err := r.DB.GetAll("databases", database.M{"tenantid": r.AuthUser.TenantId}, &databases)
	if r.Handle(err) {
		return
	}

	r.Render(http.StatusOK, databases)
}

func (c *DatabaseController) getDatabaseByShortName(r *services.ApiContext, params martini.Params) {
	var db models.Database
	err := r.DB.GetOne("databases", database.M{"tenantid": r.AuthUser.TenantId, "shortname": params["shortname"]}, &db)
	if r.Handle(err) {
		return
	}

	r.Render(http.StatusOK, db)
}

func (c *DatabaseController) createDatabase(database models.Database, r *services.ApiContext) {
	database.Init()
	database.TenantId = r.AuthUser.TenantId

	// Create backend database name
	database.Backend = fmt.Sprintf("cmdb-%s-%s", r.DB.IdToString(r.AuthUser.TenantId), database.ShortName)

	// Create database entry
	err := r.DB.Insert("databases", &database)
	if err != nil {
		log.Panic(err)
	}

	// Create actual database
	err = r.DB.CreateDatabase(database.Backend)

	r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/databases/%s", database.ShortName))
	r.Render(http.StatusCreated, "")
}
