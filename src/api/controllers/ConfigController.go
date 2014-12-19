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
	"alexandria/api/database"
	"alexandria/api/models"
	"alexandria/api/services"

	"net/http"

	"github.com/go-martini/martini"
)

type ConfigController struct {
	controller
}

func (c *ConfigController) Init(r martini.Router) error {
	// Add routes
	r.Get("/config", c.getConfig)
	return nil
}

func (c *ConfigController) getConfig(dbdriver database.Driver, r *services.ApiContext) {
	var config models.Config
	err := dbdriver.GetOne("config", nil, &config)
	r.Handle(err)

	r.Render(http.StatusOK, config)
}
