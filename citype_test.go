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
	"testing"
)

const (
	ciTypeDB = "temp"
)

func TestUpdateCIType(t *testing.T) {
	// Create a CIType
	uri := V1Uri("/cmdbs/temp/citypes")
	body := `{"name":"Original CI Type"}`
	location := Post(t, uri, body)
	defer DeleteMissing(t, location)

	// Update it
	body = `{"name":"Updated CI Name"}`
	newLocation := PutRelocate(t, location, body)
	defer Delete(t, newLocation)
}

func TestBadUpdateCIType(t *testing.T) {
	// Create a CIType
	uri := V1Uri("/cmdbs/temp/citypes")
	body := `{"name":"Original CI Type"}`
	location := Post(t, uri, body)
	defer Delete(t, location)

	// bad request
	body = `{"noname":"should fail"}`
	PutInvalid(t, location, body)
}

func TestGetAllCITypes(t *testing.T) {
	Get(t, V1Uri("/cmdbs/temp/citypes"))
}

func TestInvalidAttributeType(t *testing.T) {
	uri := V1Uri("/cmdbs/temp/citypes")

	// Test POST .../citypes with invalid attribute type
	body := `{
		"name":"BadAttributeName",
		"attributes":[
			{
				"name":"FirstAttribute",
				"type":"some_bad_type"
			}
		]}`
	PostInvalid(t, uri, body)

}

func TestNongroupWithChildren(t *testing.T) {
	uri := V1Uri("/cmdbs/temp/citypes")

	// Test POST .../citypes with invalid group attribute
	body := `{
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
}

func TestCrudCITypes(t *testing.T) {
	// Test POST .../citypes
	uri := V1Uri("/cmdbs/temp/citypes")
	body := `{
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
}
