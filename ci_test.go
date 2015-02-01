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
	body := LoadTestFixture("citype-test.json")
	typUrl := Post(t, uri, body)
	defer Delete(t, typUrl)

	// Test POST .../CI
	uri = V1Uri(fmt.Sprintf("/cmdbs/temp/%s", ciType))
	body = LoadTestFixture("ci-test.json")
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

	// test bad timestamp
	body = `{"timestamp":"The day after tomorrow", "required":false}`
	PostInvalid(t, uri, body)

	// test missing required value
	body = `{"alphanumeric":"abc123","number":"123"}`
	PostInvalid(t, uri, body)
}
