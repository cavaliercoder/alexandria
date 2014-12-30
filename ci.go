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
)

type CI struct {
	model `json:"_" bson:",inline"`

	Value map[string]interface{}
}

func (c *CI) Validate() error {
	if len(c.Value) == 0 {
		return errors.New("CI must have a valid Value body")
	}

	return nil
}

func AddCI(res http.ResponseWriter, req *http.Request) {
	// Parse request into CIType
	var ci CI
	err := Bind(req, &ci.Value)
	if Handle(res, req, err) {
		return
	}
	ci.InitModel()

	// Validate
	// TODO: Validate CI against CI Type schema
	err = ci.Validate()
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}

	// Get CMDB details
	cmdb := GetPathVar(req, "cmdb")
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		ErrNotFound(res, req)
		return
	}

	// Insert new CI
	citype := GetPathVar(req, "citype")
	err = db.C(citype).Insert(&ci)
	if Handle(res, req, err) {
		return
	}

	RenderCreated(res, req, V1Uri(fmt.Sprintf("/cmdbs/%s/%s/%s", cmdb, citype, IdToString(ci.Id))))
}

func GetCIs(res http.ResponseWriter, req *http.Request) {
	// Get CMDB details
	cmdb := GetPathVar(req, "cmdb")
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		ErrNotFound(res, req)
		return
	}

	citype := GetPathVar(req, "citype")
	var cis []CI
	err := db.C(citype).Find(nil).All(&cis)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusOK, cis)
}

func GetCIById(res http.ResponseWriter, req *http.Request) {
	// Get CMDB details
	cmdb := GetPathVar(req, "cmdb")
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		ErrNotFound(res, req)
		return
	}

	// Get collection
	citype := GetPathVar(req, "citype")

	// Get Id
	id := GetPathVar(req, "id")
	oid, err := IdFromString(id)
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}

	var ci CI
	err = db.C(citype).FindId(oid).One(&ci)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusOK, ci)
}

func DeleteCIById(res http.ResponseWriter, req *http.Request) {
	cmdb := GetPathVar(req, "cmdb")
	citype := GetPathVar(req, "citype")
	id := GetPathVar(req, "id")

	// Get CMDB details
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		ErrNotFound(res, req)
		return
	}

	// Get id
	oid, err := IdFromString(id)
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}

	// Remove the CI
	err = db.C(citype).RemoveId(oid)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusNoContent, "")
}
