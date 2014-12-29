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
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func get(t *testing.T, uri string, code int) {
	fmt.Printf("[TEST] GET %s (expecting %d)...\n", uri, code)

	// Create request
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	expect(t, res.Code, code)
}

func Get(t *testing.T, uri string) {
	get(t, uri, http.StatusOK)
}

func GetMissing(t *testing.T, uri string) {
	get(t, uri, http.StatusNotFound)
}

func post(t *testing.T, uri string, body string, code int) string {
	fmt.Printf("[TEST] POST %s (expecting %d)...\n", uri, code)

	// Create request
	reader := strings.NewReader(body)
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		panic(err)
	}

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	expect(t, res.Code, code)

	return res.HeaderMap.Get("Location")
}

func Post(t *testing.T, uri string, body string, testDuplicates bool) {
	// Create a new resource
	location := post(t, uri, body, http.StatusCreated)

	if location == "" {
		t.Errorf("No location header was returned for new resource in: %s", uri)
	} else {
		// Retrieve the new resource
		Get(t, location)

		if testDuplicates {
			post(t, uri, body, http.StatusConflict)
		}

		// Delete the resource
		Delete(t, location)

		// Ensure it is missing
		GetMissing(t, location)
	}
}

func Delete(t *testing.T, uri string) {
	fmt.Printf("[TEST] DELETE %s...\n", uri)

	// Create request
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		panic(err)
	}

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	expect(t, res.Code, http.StatusNoContent)
}
