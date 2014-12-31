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
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	hash := HashPassword("Password1")

	if hashlen := len(hash); hashlen < 64 {
		t.Errorf("Expected password hash to be at least 64 characters but it was only %d: %s", hashlen, hash)
	}

	if !CheckPassword(hash, "Password1") {
		t.Error("Expected correct password to validate but it did not")
	}

	if CheckPassword(hash, "WrongP4ssw0RD!") {
		t.Error("Expected incorrect password to fail validation but it passed")
	}
}
