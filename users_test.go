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
	testEmail     = "test.user@cavaliercoder.com"
	testFirstName = "Test"
	testLastName  = "User"
)

func TestAddUser(t *testing.T) {
	// Test POST /users
	uri := fmt.Sprintf("%s/users", ApiV1Prefix)
	body := fmt.Sprintf(`{"email":"%s","firstName":"%s","lastName":"%s"}`, testEmail, testFirstName, testLastName)
	Post(t, uri, body, true)
}

func TestGetUsers(t *testing.T) {
	// Test GET /users
	Get(t, fmt.Sprintf("%s/users", ApiV1Prefix))
}
