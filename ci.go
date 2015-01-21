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
 */
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type CI struct {
	model `json:"-" xml:"-" bson:",inline"`

	Value map[string]interface{}
}

func (c *CI) Validate() error {
	if len(c.Value) == 0 {
		return errors.New("CI must have a valid Value body")
	}

	return nil
}

func AddCI(res http.ResponseWriter, req *http.Request) {
	cmdb := GetPathVar(req, "cmdb")
	citype := GetPathVar(req, "citype")

	// Parse request into CIType
	var ci CI
	err := Bind(req, &ci.Value)
	if Handle(res, req, err) {
		return
	}
	ci.InitModel()

	// Get CMDB details
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		log.Printf("No such CMDB found: %s", cmdb)
		ErrNotFound(res, req)
		return
	}

	// Get CI Type schema
	var typ CIType
	err = db.C("citypes").Find(M{"shortname": citype}).One(&typ)
	if Handle(res, req, err) {
		log.Printf("No such CI type found: %s", citype)
		return
	}

	// Validate parser
	err = ci.Validate()
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}

	// Validate against schema
	err = validateFields(&ci.Value, &typ.Attributes, "")
	if err != nil {
		ErrBadRequest(res, req, err)
		return
	}

	// Insert new CI
	err = db.C(citype).Insert(&ci)
	if Handle(res, req, err) {
		return
	}

	RenderCreated(res, req, V1Uri(fmt.Sprintf("/cmdbs/%s/%s/%s", cmdb, citype, IdToString(ci.Id))))
}

func validateFields(fields *map[string]interface{}, schema *CITypeAttributeList, path string) error {
	for key, _ := range *fields {
		fullPath := fmt.Sprintf("%s.%s", path, key)

		// Dereference the value so it may be modified by format.Validate()
		val := (*fields)[key]

		// Does this key exist in the schema?
		att := schema.Get(key)
		if att == nil {
			return errors.New(fmt.Sprintf("No schema definition found for field '%s'", fullPath))
		}

		// Does the format exist?
		format := GetAttributeFormat(att.Type)
		if format == nil {
			return errors.New(fmt.Sprintf("No format parser found for type '%s' in field '%s'", att.Type, fullPath))
		}

		// Is the value valid?
		// This will also translate the value if required
		// TODO: Need to dereference val in the range loop so the original object is updated
		err := format.Validate(att, &val)
		if err != nil {
			return err
		}

		// Process children?
		if len(att.Children) > 0 {
			childFields, ok := val.(map[string]interface{})
			if !ok {
				return errors.New(fmt.Sprintf("Expected '%s' to be a valid JSON object", fullPath))
			}

			err = validateFields(&childFields, &att.Children, fullPath)
			if err != nil {
				return err
			}
		}
	}

	// Ensure all required fields were included
	for _, att := range *schema {
		if att.Required {
			if _, ok := (*fields)[att.ShortName]; !ok {
				return errors.New(fmt.Sprintf("Required field '%s' is not present", att.Name))
			}
		}
	}

	return nil
}

func GetCIs(res http.ResponseWriter, req *http.Request) {
	// Get CMDB details
	cmdb := GetPathVar(req, "cmdb")
	db := GetCmdbBackend(req, cmdb)
	if db == nil {
		log.Printf("No such CMDB found: %s", cmdb)
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
		log.Printf("No such CMDB found: %s", cmdb)
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
		log.Printf("No such CMDB found: %s", cmdb)
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
