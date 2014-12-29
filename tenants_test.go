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
	"testing"
)

const (
	testTenant = "Test tenant"
)

func TestAddTenant(t *testing.T) {
	n := GetServer()
	res := httptest.NewRecorder()

	// Test POST /users
	reqBody := fmt.Sprintf(`{"name":"%s"}`, testTenant)
	req := Post("/tenants", reqBody)
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusCreated)

	// Test returned location header
	location := res.HeaderMap.Get("Location")
	if location == "" {
		t.Errorf("No location header was set for a created tenant resource:\n%#v", res.HeaderMap)
	} else {
		// Test GET tenant by code
		res = httptest.NewRecorder()
		req = Get(location)
		n.ServeHTTP(res, req)
		expect(t, res.Code, http.StatusOK)

		// Test DELETE tenant by code
		res = httptest.NewRecorder()
		req = Delete(location)
		n.ServeHTTP(res, req)
		expect(t, res.Code, http.StatusNoContent)

		// Test GET missing tenant
		res = httptest.NewRecorder()
		req = Get(location)
		n.ServeHTTP(res, req)
		expect(t, res.Code, http.StatusNotFound)
	}
}

func TestGetTenants(t *testing.T) {
	n := GetServer()
	res := httptest.NewRecorder()

	// Test GET /users
	req := Get("/tenants")
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusOK)
}
