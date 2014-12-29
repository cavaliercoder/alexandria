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
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const (
	testEmail     = "test.user@cavaliercoder.com"
	testFirstName = "Test"
	testLastName  = "User"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Get(uri string) *http.Request {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		panic(err)
	}

	return req
}

func Post(uri string, body string) *http.Request {
	var reader io.Reader = nil

	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		panic(err)
	}

	return req
}

func Delete(uri string) *http.Request {
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		panic(err)
	}

	return req
}

func TestAddUser(t *testing.T) {
	n := GetServer()
	res := httptest.NewRecorder()

	// Test POST /users
	reqBody := fmt.Sprintf(`{"email":"%s","firstName":"%s","lastName":"%s"}`, testEmail, testFirstName, testLastName)
	req := Post("/users", reqBody)
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusCreated)

	// Make sure duplicates cant be created
	req = Post("/users", reqBody)
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusConflict)
}

func TestGetUsers(t *testing.T) {
	n := GetServer()
	res := httptest.NewRecorder()

	// Test GET /users
	req := Get("/users")
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusOK)

	// Test GET /users/:email
	req = Get(fmt.Sprintf("/users/%s", testEmail))
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusOK)

	// Test GET /users/:missing
	res = httptest.NewRecorder()
	req = Get("/users/i_dont_exist")
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusNotFound)
}

func TestDeleteUser(t *testing.T) {
	n := GetServer()
	res := httptest.NewRecorder()

	// Test POST /users
	req := Delete(fmt.Sprintf("/users/%s", testEmail))
	n.ServeHTTP(res, req)
	expect(t, res.Code, http.StatusNoContent)
}
