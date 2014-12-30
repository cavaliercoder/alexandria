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
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type Tenant struct {
	model `json:"-" bson:",inline"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Cmdbs map[string]Cmdb `json:"-"`
}

func (c *Tenant) InitModel() {
	c.model.InitModel()

	// Create the code by hashing the database ID
	hash := sha256.Sum256([]byte(IdToString(c.Id)))
	c.Code = fmt.Sprintf("%x-%x-%x", hash[0:2], hash[3:6], hash[7:10])
}

func (c *Tenant) Validate() error {
	if c.Code == "" {
		return errors.New("No tenancy code specified")
	}

	matched, _ := regexp.MatchString("^[a-f0-9]{4}-[a-f0-9]{6}-[a-f0-9]{6}$", c.Code)
	if !matched {
		return errors.New(fmt.Sprintf("Invalid tenancy code: %s", c.Code))
	}

	if c.Name == "" {
		return errors.New("No tenant name specified")
	}

	return nil
}

func GetTenants(res http.ResponseWriter, req *http.Request) {
	var tenants []Tenant
	err := RootDb().C("tenants").Find(nil).All(&tenants)
	Handle(res, req, err)

	Render(res, req, http.StatusOK, tenants)
}

func GetTenantByCode(res http.ResponseWriter, req *http.Request) {
	code := GetPathVar(req, "code")

	var tenant Tenant
	err := RootDb().C("tenants").Find(M{"code": code}).One(&tenant)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusOK, tenant)
}

func AddTenant(res http.ResponseWriter, req *http.Request) {
	var tenant Tenant
	err := Bind(req, &tenant)
	if Handle(res, req, err) {
		return
	}
	tenant.InitModel()

	err = tenant.Validate()
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}

	err = RootDb().C("tenants").Insert(&tenant)
	if Handle(res, req, err) {
		return
	}

	RenderCreated(res, req, V1Uri(fmt.Sprintf("/tenants/%s", tenant.Code)))
}

func DeleteTenantByCode(res http.ResponseWriter, req *http.Request) {
	code := GetPathVar(req, "code")

	// TODO: Ensure only users for current tenant can be deleted
	err := RootDb().C("tenants").Remove(M{"code": code})
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusNoContent, "")
}
