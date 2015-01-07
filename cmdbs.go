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
 */
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Cmdb struct {
	model       `json:"-" bson:",inline"`
	TenantId    interface{} `json:"-"`
	Name        string      `json:"name"`
	ShortName   string      `json:"shortName"`
	Description string      `json:"description"`
}

func (c *Cmdb) Validate() error {
	if c.Name == "" {
		return errors.New("No CMDB name specified")
	}

	if c.ShortName == "" {
		c.ShortName = GetShortName(c.Name)
	}

	if !IsValidShortName(c.ShortName) {
		return errors.New(fmt.Sprintf("Invalid characters in CMDB name: '%s'", c.ShortName))
	}

	if c.TenantId == nil {
		return errors.New("No tenancy code specified")
	}

	return nil
}

func (c *Cmdb) GetBackendName() string {
	return fmt.Sprintf("cmdb_%s", IdToString(c.Id))
}

func GetCmdbs(res http.ResponseWriter, req *http.Request) {
	auth := GetAuthContext(req)

	v := make([]Cmdb, 0, len(auth.Tenant.Cmdbs))
	for _, value := range auth.Tenant.Cmdbs {
		v = append(v, value)
	}

	Render(res, req, http.StatusOK, v)
}

func GetCmdbByName(res http.ResponseWriter, req *http.Request) {
	auth := GetAuthContext(req)
	name := GetPathVar(req, "name")

	cmdb, ok := auth.Tenant.Cmdbs[name]
	if !ok {
		ErrNotFound(res, req)
		return
	}

	Render(res, req, http.StatusOK, cmdb)
}

func AddCmdb(res http.ResponseWriter, req *http.Request) {
	auth := GetAuthContext(req)

	// Parse request and bind to Cmdb{}
	var cmdb Cmdb
	err := Bind(req, &cmdb)
	if Handle(res, req, err) {
		return
	}

	cmdb.InitModel()
	cmdb.TenantId = auth.User.TenantId

	// Validate
	err = cmdb.Validate()
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}

	// Prevent duplicates
	_, ok := auth.Tenant.Cmdbs[cmdb.ShortName]
	if ok {
		log.Printf("Bad request: A CMDB already exists with name '%s'", cmdb.ShortName)
		ErrConflict(res, req)
		return
	}

	// Insert in database
	field := fmt.Sprintf("cmdbs.%s", cmdb.ShortName)
	mgoErr := RootDb().C("tenants").Update(M{"_id": auth.User.TenantId}, M{"$set": M{field: &cmdb}})
	if Handle(res, req, mgoErr) {
		return
	}

	// Create backend
	err = CreateCmdb(cmdb.GetBackendName())
	if Handle(res, req, err) {
		return
	}

	// Tell the world
	RenderCreated(res, req, V1Uri(fmt.Sprintf("/cmdbs/%s", cmdb.ShortName)))
}

func DeleteCmdbByName(res http.ResponseWriter, req *http.Request) {
	auth := GetAuthContext(req)
	name := GetPathVar(req, "name")

	cmdb, ok := auth.Tenant.Cmdbs[name]
	if !ok {
		ErrNotFound(res, req)
		return
	}

	field := fmt.Sprintf("cmdbs.%s", cmdb.ShortName)
	mgoErr := RootDb().C("tenants").Update(M{"_id": auth.User.TenantId}, M{"$unset": M{field: ""}})
	if Handle(res, req, mgoErr) {
		return
	}

	// Drop backend
	err := DropCmdb(cmdb.GetBackendName())
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusNoContent, "")
}
