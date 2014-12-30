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
	"fmt"
	"testing"
)

const (
	ciTypeDB          = "CITypeDB"
	ciTypeName        = "TestCIType"
	ciTypeDescription = "Test CI Type"
)

func TestCITypes(t *testing.T) {
	// Create a cmdb
	uri := V1Uri("/cmdbs")
	body := fmt.Sprintf(`{"name":"%s"}`, ciTypeDB)
	dburl := Post(t, uri, body)
	defer Delete(t, dburl)

	// Test POST /cmdbs
	uri = V1Uri(fmt.Sprintf("/cmdbs/%s/citypes", ciTypeDB))
	body = fmt.Sprintf(`{"name":"%s","description":"%s"}`, ciTypeName, ciTypeDescription)
	Crud(t, uri, body, true)

	body = `{"name":"Invalid Name!"}`
	PostInvalid(t, uri, body)

	// Test get all
	Get(t, uri)
}
