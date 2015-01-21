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
	"fmt"
	"testing"
)

const (
	ciType = "test-ci-type"
)

func TestCI(t *testing.T) {
	// Create temporary CI Type
	uri := V1Uri("/cmdbs/temp/citypes")
	body := `{
		"name":"Test CI Type",
		"description": "A test CI Type",
		"attributes": [
			{
				"name":"alphanumeric",
				"type":"string",
				"filters": ["^[A-Za-z0-9]+$"]
			},
			{
				"name":"number",
				"type":"number",
				"minValue":100,
				"maxValue":200
			},
			{
				"name":"Required",
				"type":"boolean",
				"required":true
			},
			{
				"name":"group",
				"type":"group",
				"children":[
					{
						"name":"allCaps",
						"type":"string",
						"filters":["^[A-Z]+$"]
					},
					{
						"name":"grandchildren",
						"type":"group",
						"children":[
							{
								"name":"grandchild",
								"type":"string"
							}
						]
					}
				]
			}
		]
		}`
	typUrl := Post(t, uri, body)
	defer Delete(t, typUrl)

	// Test POST .../CI
	uri = V1Uri(fmt.Sprintf("/cmdbs/temp/%s", ciType))
	body = `{
		"alphanumeric":"StringValue123",
		"number":123,
		"required":"Yes",
		"group":{
			"allCaps":"ABC",
			"grandchildren":{
				"grandchild":"Some value"
			}
		}
	}`
	Crud(t, uri, body, false)

	// test POST invalid CI schema
	body = `{"badAttribute":"some value", "required":false}`
	PostInvalid(t, uri, body)

	// test POST string filter
	body = `{"alphanumeric":"Not @lphANUM3r1c!", "required": false}`
	PostInvalid(t, uri, body)

	// test number minimum
	body = `{"number":1, "required":false}`
	PostInvalid(t, uri, body)

	// test number maximum
	body = `{"number":321, "required":false}`
	PostInvalid(t, uri, body)

	// test missing required value
	body = `{"alphanumeric":"abc123","number":"123"}`
	PostInvalid(t, uri, body)
}
