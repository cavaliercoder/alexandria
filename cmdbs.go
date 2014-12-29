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
package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type Cmdb struct {
	model       `json:"-" bson:",inline"`
	TenantId    interface{} `json:"-"`
	Name        string      `json:"name"`
	ShortName   string      `json:"shortName"`
	Description string      `json:"description"`
}

func (c *Cmdb) Validate() error {
	if match, _ := regexp.MatchString("^[a-zA-Z0-9-_]+$", c.ShortName); !match {
		return errors.New("Invalid short name")
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
	authUser := GetApiUser(req)

	var cmdbs []Cmdb
	err := RootDb().C("cmdbs").Find(M{"tenantid": authUser.TenantId}).All(&cmdbs)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusOK, cmdbs)
}

func GetCmdbByName(res http.ResponseWriter, req *http.Request) {
	authUser := GetApiUser(req)
	name := GetPathVar(req, "name")

	var cmdb Cmdb
	err := RootDb().C("cmdbs").Find(M{"tenantid": authUser.TenantId, "shortname": name}).One(&cmdb)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusOK, cmdb)
}

func AddCmdb(res http.ResponseWriter, req *http.Request) {
	authUser := GetApiUser(req)
	var cmdb Cmdb
	err := Bind(req, &cmdb)
	if Handle(res, req, err) {
		return
	}

	cmdb.InitModel()
	cmdb.TenantId = authUser.TenantId

	// Validate
	err = cmdb.Validate()
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}

	err = RootDb().C("cmdbs").Insert(&cmdb)
	if Handle(res, req, err) {
		return
	}

	RenderCreated(res, req, fmt.Sprintf("%s/cmdbs/%s", ApiV1Prefix, cmdb.ShortName))
}

func DeleteCmdbByName(res http.ResponseWriter, req *http.Request) {
	authUser := GetApiUser(req)
	name := GetPathVar(req, "name")

	err := RootDb().C("cmdbs").Remove(M{"tenantid": authUser.TenantId, "shortname": name})
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusNoContent, "")
}
