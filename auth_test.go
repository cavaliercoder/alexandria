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
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	// A URI that requires X-AUTH-TOKEN header
	authTestUri = "/cmdbs"
)

func TestNoAuthHeader(t *testing.T) {
	code := http.StatusUnauthorized
	uri := V1Uri(authTestUri)

	fmt.Printf("[TEST] GET %s (expecting %d)...\n", uri, code)

	// Create request
	req, _ := http.NewRequest("GET", uri, nil)

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	areEqual(t, res.Code, code)
}

func TestBadAuthHeader(t *testing.T) {
	code := http.StatusUnauthorized
	uri := V1Uri(authTestUri)

	fmt.Printf("[TEST] GET %s (expecting %d)...\n", uri, code)

	// Create request
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("X-Auth-Token", "abc123")

	// Create response recorder
	res := httptest.NewRecorder()

	// Start web server
	n := GetServer()
	n.ServeHTTP(res, req)

	// Validate response
	areEqual(t, res.Code, code)
}

func TestGoodAuthHeader(t *testing.T) {
	Get(t, V1Uri(authTestUri))
}
