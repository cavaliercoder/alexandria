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
	"strings"
)

const (
	ciTypeCollection = "citypes"
)

type CIType struct {
	model `json:"-" bson:",inline"`

	InheritFrom string              `json:"inheritFrom"`
	Name        string              `json:"name"`
	ShortName   string              `json:"shortName"`
	Description string              `json:"description"`
	Attributes  CITypeAttributeList `json:"attributes"`
}

type CITypeAttribute struct {
	Name        string              `json:"name"`
	ShortName   string              `json:"shortName"`
	Description string              `json:"description"`
	Type        string              `json:"type"`
	MinValues   int                 `json:"minValues"`
	MaxValues   int                 `json:"maxValues"`
	Filter      string              `json:"filter"`
	Children    CITypeAttributeList `json:"children"`
}

type CITypeAttributeList []CITypeAttribute

func (c *CITypeAttributeList) Get(name string) *CITypeAttribute {
	name = strings.ToLower(name)
	for _, att := range *c {
		if att.ShortName == name {
			return &att
		}
	}

	return nil
}

func (c *CIType) Validate() error {
	if c.Name == "" {
		return errors.New("No CI Type name specified")
	}

	if c.ShortName == "" {
		c.ShortName = GetShortName(c.Name)
	}

	if !IsValidShortName(c.ShortName) {
		return errors.New("Invalid characters in CI Type name")
	}

	// Validate each attribute
	err := c.validateAttributes(&c.Attributes, "")
	if err != nil {
		return err
	}

	return nil
}

func (c *CIType) validateAttributes(atts *CITypeAttributeList, path string) error {
	for index, _ := range *atts {
		// Derefence the attribute so it may be modified
		att := &(*atts)[index]

		if att.Name == "" {
			return errors.New("No attribute name specified")
		}

		if att.ShortName == "" {
			att.ShortName = GetShortName(att.Name)
		}

		if !IsValidShortName(att.ShortName) {
			return errors.New(fmt.Sprintf("Invalid characters in CI Attribute '%s%s'", path, att.ShortName))
		}

		if att.Type == "" {
			return errors.New(fmt.Sprintf("No type specified for CI Attribute '%s%s'", path, att.ShortName))
		}

		if GetAttributeFormat(att.Type) == nil {
			return errors.New(fmt.Sprintf("Unsupported attribute format '%s' for CI Attribute '%s%s'", att.Type, path, att.ShortName))
		}

		if att.Type == "group" {
			// Validate children
			err := c.validateAttributes(&att.Children, fmt.Sprintf("%s.", att.ShortName))
			if err != nil {
				return err
			}
		} else if len(att.Children) > 0 {
			return errors.New(fmt.Sprintf("CI Attribute '%s%s' has children but is not a group attribute", path, att.ShortName))
		}
	}

	return nil
}

func GetCITypes(res http.ResponseWriter, req *http.Request) {
	// Get CMDB details
	cmdb := GetPathVar(req, "cmdb")
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		ErrNotFound(res, req)
		return
	}

	var citypes []CIType
	err := db.C(ciTypeCollection).Find(nil).All(&citypes)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusOK, citypes)
}

func GetCITypeByName(res http.ResponseWriter, req *http.Request) {
	// Get CMDB details
	cmdb := GetPathVar(req, "cmdb")
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		ErrNotFound(res, req)
		return
	}

	// Get the type
	var citype CIType
	name := GetPathVar(req, "name")
	err := db.C(ciTypeCollection).Find(M{"shortname": name}).One(&citype)
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusOK, citype)
}

func AddCIType(res http.ResponseWriter, req *http.Request) {
	// Parse request into CIType
	var citype CIType
	err := Bind(req, &citype)
	if Handle(res, req, err) {
		return
	}
	citype.InitModel()

	// Validate
	err = citype.Validate()
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

	// Insert new type
	err = db.C(ciTypeCollection).Insert(&citype)
	if Handle(res, req, err) {
		return
	}

	RenderCreated(res, req, V1Uri(fmt.Sprintf("/cmdbs/%s/citypes/%s", cmdb, citype.ShortName)))
}

func DeleteCITypeByName(res http.ResponseWriter, req *http.Request) {
	cmdb := GetPathVar(req, "cmdb")
	name := GetPathVar(req, "name")

	// Get CMDB details
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		ErrNotFound(res, req)
		return
	}

	err := db.C(ciTypeCollection).Remove(M{"shortname": name})
	if Handle(res, req, err) {
		return
	}

	Render(res, req, http.StatusNoContent, "")
}
