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

func TestGetUsers(t *testing.T) {
	n := GetServer()
	res := httptest.NewRecorder()

	// Test GET /users
	req, _ := http.NewRequest("GET", "/users", nil)
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusOK)
}

func TestAddUser(t *testing.T) {
	n := GetServer()
	res := httptest.NewRecorder()

	// Test POST /users
	reqBody := `{"email":"testuser@alexandria.org"}`
	reader := strings.NewReader(reqBody)
	req, _ := http.NewRequest("POST", "/users", reader)
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusCreated)
}

func TestDeleteUser(t *testing.T) {
	n := GetServer()
	res := httptest.NewRecorder()

	// Test POST /users
	req, _ := http.NewRequest("DELETE", "/users/testuser@alexandria.org", nil)
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusNoContent)
}
