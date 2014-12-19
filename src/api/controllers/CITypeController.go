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

const (
	citype_table string = "_citypes"
)

type CITypeController struct {
	controller
}

func (c *CITypeController) Init(r martini.Router) error {

	// Add routes
	r.Get("/citypes", c.getCITypes)
	r.Get("/citypes/:shortname", c.getCITypeByShortName)
	r.Post("/citypes", binding.Bind(models.CIType{}), c.addCIType)
        r.Delete("/citypes/:shortname", c.deleteCITypeByShortName)

	return nil
}

func (c *CITypeController) getCITypes(r *services.ApiContext) {
	var citypes []models.CIType
	err := r.DB.GetAll(citype_table, nil, &citypes)
	r.Handle(err)

	r.Render(http.StatusOK, citypes)
}

func (c *CITypeController) getCITypeByShortName(r *services.ApiContext, params martini.Params) {
	var citype models.CIType
	err := r.DB.GetOne(citype_table, database.M{"shortname": params["shortname"]}, &citype)
	if r.Handle(err) { return }

	r.Render(http.StatusOK, citype)
}

func (c *CITypeController) addCIType(citype models.CIType, r *services.ApiContext) {
	citype.Init()
	citype.TenantId = r.AuthUser.TenantId
        
        err := r.DB.Insert(citype_table, citype)
	if err != nil { log.Panic(err) }

	r.ResponseWriter.Header().Set("Location", fmt.Sprintf("/citypes/%s", citype.ShortName))
	r.Render(http.StatusCreated, "")
}

func (c *CITypeController) deleteCITypeByShortName(r *services.ApiContext, params martini.Params) {
        err := r.DB.Delete(citype_table, database.M{"tenantid": r.AuthUser.TenantId, "shortname": params["shortname"]})
        if r.Handle(err) { return }
        
        r.Render(http.StatusNoContent, "")
}
