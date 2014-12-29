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
	testTenant = "Test tenant"
)

func TestAddTenant(t *testing.T) {
	// Test POST /users
	reqBody := fmt.Sprintf(`{"name":"%s"}`, testTenant)
	Post(t, fmt.Sprintf("%s/tenants", ApiV1Prefix), reqBody, false)
}

func TestGetTenants(t *testing.T) {
	// Test GET /users
	Get(t, fmt.Sprintf("%s/tenants", ApiV1Prefix))
}
