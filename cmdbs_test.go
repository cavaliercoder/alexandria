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
	cmdbName        = "TestCMDB with 0dd n@me!"
	cmdbDescription = "A temporary test CMDB"
)

func TestAddCmdb(t *testing.T) {
	// Test POST /cmdbs
	uri := V1Uri("/cmdbs")
	body := fmt.Sprintf(`{"name":"%s","description":"%s"}`, cmdbName, cmdbDescription)
	Crud(t, uri, body, true)
}

func TestGetCmdbs(t *testing.T) {
	// Test GET /cmdbs
	Get(t, V1Uri("/cmdbs"))
}
