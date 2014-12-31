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
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func areEqual(t *testing.T, a interface{}, b interface{}) bool {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
		return false
	}

	return true
}

func getRootUser() *User {
	// Get apiInfo
	var apiInfo ApiInfo
	err := RootDb().C("apiInfo").Find(nil).One(&apiInfo)
	if err != nil {
		return nil
	}

	var user User
	err = RootDb().C("users").FindId(apiInfo.RootUserId).One(&user)
	if err != nil {
		return nil
	}

	return &user
}

func NewRequest(method string, uri string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		panic(err)
	}

	user := getRootUser()

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("X-Auth-Token", user.ApiKey)
	req.Header.Add("User-Agent", "Alexandria CMDB Tests")

	return req
}

func post(t *testing.T, uri string, body string, code int) string {
	fmt.Printf("[TEST] POST %s (expecting %d)...\n", uri, code)

	// Create request
	reader := strings.NewReader(body)
	req := NewRequest("POST", uri, reader)

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	areEqual(t, res.Code, code)

	return res.HeaderMap.Get("Location")
}

// Post posts the specified resource and expects a 201 Created response
func Post(t *testing.T, uri string, body string) string {
	return post(t, uri, body, http.StatusCreated)
}

// Post posts the specified invalid resource and expects a 400 Bad request response
func PostInvalid(t *testing.T, uri string, body string) {
	post(t, uri, body, http.StatusBadRequest)
}

func patch(t *testing.T, uri string, body string, code int) {
	fmt.Printf("[TEST] PATCH %s (expecting %d)...\n", uri, code)

	// Create request
	reader := strings.NewReader(body)
	req := NewRequest("PATCH", uri, reader)

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	areEqual(t, res.Code, code)
}

// Patch posts the specified resource update and expects a 204 No content repsonse
func Patch(t *testing.T, uri string, body string) {
	patch(t, uri, body, http.StatusNoContent)
}

// Patch posts the specified invalid resource patch and expects a 400 Bad request response
func PatchInvalid(t *testing.T, uri string, body string) {
	patch(t, uri, body, http.StatusBadRequest)
}

func get(t *testing.T, uri string, code int) {
	fmt.Printf("[TEST] GET %s (expecting %d)...\n", uri, code)

	// Create request
	req := NewRequest("GET", uri, nil)

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	areEqual(t, res.Code, code)
}

// Get retrieves a resource and expects a 200 Ok response
func Get(t *testing.T, uri string) {
	get(t, uri, http.StatusOK)
}

// Get retrieves a nonexistant resource and expects a 404 Not found response
func GetMissing(t *testing.T, uri string) {
	get(t, uri, http.StatusNotFound)
}

// Delete deletes a resource and expects a 204 No content response
func Delete(t *testing.T, uri string) {
	fmt.Printf("[TEST] DELETE %s...\n", uri)

	// Create request
	req := NewRequest("DELETE", uri, nil)

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	areEqual(t, res.Code, http.StatusNoContent)
}

// Crud tests the creation, retrieval, update and deletion of an api resource.
// If testDuplicates is true, Crud attempts to create a duplicate resource and
// expects a 409 Conflict response.
func Crud(t *testing.T, uri string, body string, testDuplicates bool) {
	// Create a new resource
	location := Post(t, uri, body)

	if location == "" {
		t.Errorf("No location header was returned for new resource in: %s", uri)
	} else {
		// Retrieve the new resource
		Get(t, location)

		// Make sure duplicates can't be created
		if testDuplicates {
			post(t, uri, body, http.StatusConflict)
		}

		// Delete the resource
		Delete(t, location)

		// Ensure it is missing
		GetMissing(t, location)
	}
}
