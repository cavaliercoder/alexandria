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
	"fmt"
	"testing"
)

const (
	ciTypeDB = "CITypeDB"
)

func TestUpdateCIType(t *testing.T) {
	// Create a temporary cmdb
	uri := V1Uri("/cmdbs")
	body := fmt.Sprintf(`{"name":"%s"}`, ciTypeDB)
	dburl := Post(t, uri, body)
	defer Delete(t, dburl)

	// Create a CIType
	uri = V1Uri(fmt.Sprintf("/cmdbs/%s/citypes", ciTypeDB))
	body = `{"name":"Original CI Type"}`
	location := Post(t, uri, body)
	defer DeleteMissing(t, location)

	// Update it
	body = `{"name":"Updated CI Name"}`
	newLocation := PutRelocate(t, location, body)
	defer Delete(t, newLocation)

	// bad request
	body = `{"noname":"should fail"}`
	PutInvalid(t, newLocation, body)
}

func TestCITypes(t *testing.T) {
	// Create a temporary cmdb
	uri := V1Uri("/cmdbs")
	body := fmt.Sprintf(`{"name":"%s"}`, ciTypeDB)
	dburl := Post(t, uri, body)
	defer Delete(t, dburl)

	// Test POST .../citypes
	uri = V1Uri(fmt.Sprintf("/cmdbs/%s/citypes", ciTypeDB))
	body = `{
		"name":"Test_CI_Type with w!3rd CH@RS!",
		"description": "A test CI Type",
		"attributes": [
			{
				"name":"FirstAttribute",
				"description": "The first attribute",
				"type":"string"
			},
			{
				"name":"SecondAttribute",
				"description":"The second attribute (with children)",
				"type":"group",
				"children":[
					{
						"name":"GrandchildAttribute",
						"description":"Grandchild Attribute",
						"type":"string"
					}
				]
			}
		]
		}`
	Crud(t, uri, body, true)

	// Test POST .../citypes with invalid attribute type
	body = `{
		"name":"BadAttributeName",
		"attributes":[
			{
				"name":"FirstAttribute",
				"type":"some_bad_type"
			}
		]}`
	PostInvalid(t, uri, body)

	// Test POST .../citypes with invalid group attribute
	body = `{
		"name":"BadGroupAttribute",
		"attributes":[
			{
				"name":"BadAttribute",
				"type":"string",
				"children":[
					{
						"name":"Impossible",
						"type":"string"
					}
				]
			}
		]}`
	PostInvalid(t, uri, body)

	// Test get all
	Get(t, uri)
}
