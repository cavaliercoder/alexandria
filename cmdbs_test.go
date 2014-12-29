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
	cmdbName        = "Test CMDB"
	cmdbShortName   = "test"
	cmdbDescription = "A temporary test CMDB"
)

func TestAddCmdb(t *testing.T) {
	// Test POST /cmdbs
	uri := fmt.Sprintf("%s/cmdbs", ApiV1Prefix)
	body := fmt.Sprintf(`{"name":"%s","shortName":"%s","description":"%s"}`, cmdbName, cmdbShortName, cmdbDescription)
	Post(t, uri, body, true)

	// TODO: Add a test to ensure invalid cmdb creation fails (i.e field validation)
}

func TestGetCmdbs(t *testing.T) {
	// Test GET /cmdbs
	Get(t, fmt.Sprintf("%s/cmdbs", ApiV1Prefix))
}
