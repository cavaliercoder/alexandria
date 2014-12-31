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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testEmail     = "test.user@cavaliercoder.com"
	testFirstName = "Test"
	testLastName  = "User"
	testPassword  = "Password1"
)

func TestAddUser(t *testing.T) {
	// Test POST /users
	uri := V1Uri("/users")
	body := fmt.Sprintf(`{"email":"%s","firstName":"%s","lastName":"%s","password":"%s"}`, testEmail, testFirstName, testLastName, testPassword)
	Crud(t, uri, body, true)

	// Prevent missing email addresses
	body = `{"firstName":"No","lastName":"Email","password":"Password1"}`
	PostInvalid(t, uri, body)

	// Prevent invalid email addresses
	body = `{"email":"not valid email address","password":"Password1"}`
	PostInvalid(t, uri, body)
}

func TestUserPassword(t *testing.T) {
	// Create a temporary user
	uri := V1Uri("/users")
	body := fmt.Sprintf(`{"email":"%s","firstName":"%s","lastName":"%s","password":"%s"}`, testEmail, testFirstName, testLastName, testPassword)
	userurl := Post(t, uri, body)
	defer Delete(t, userurl)

	// Update the password
	uri = fmt.Sprintf("%s/password", userurl)
	password := fmt.Sprintf(`{"password":"%s"}`, testPassword)
	Patch(t, uri, password)

	password = `{"invalid":true}`
	PatchInvalid(t, uri, password)

	// Test login
	testLogin(t, testEmail, testPassword, http.StatusOK)
	testLogin(t, testEmail, "BadPassword", http.StatusUnauthorized)
	testLogin(t, "i_dont_exist", "AnyPassword", http.StatusUnauthorized)
}

func testLogin(t *testing.T, username string, password string, code int) {
	uri := V1Uri("/apikey")
	fmt.Printf("[TEST] POST %s (expecting %d)...\n", uri, code)

	// Create request
	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	req := NewRequest("POST", uri, strings.NewReader(body))
	req.Header.Del("X-Auth-Token") // Ensure the route works without API key auth
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	areEqual(t, res.Code, code)
}

func TestGetUsers(t *testing.T) {
	// Test GET /users
	Get(t, V1Uri("/users"))

	Get(t, V1Uri("/users/current"))
}
