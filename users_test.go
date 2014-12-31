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
	uri := V1Uri("/users")
	body := fmt.Sprintf(`{"email":"%s","firstName":"%s","lastName":"%s"}`, testEmail, testFirstName, testLastName)
	Crud(t, uri, body, true)

	// Prevent missing email addresses
	body = `{"firstName":"No","lastName":"Email"}`
	PostInvalid(t, uri, body)

	// Prevent invalid email addresses
	body = `{"email":"not valid email address"}`
	PostInvalid(t, uri, body)
}

func TestGetUsers(t *testing.T) {
	// Test GET /users
	Get(t, V1Uri("/users"))

	Get(t, V1Uri("/users/current"))
}
