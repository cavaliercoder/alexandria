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

func TestStringFormat(t *testing.T) {
	format := GetAttributeFormat("string")
	if format == nil {
		t.Errorf("String attribute format does not appear to be registered")
		return
	}

	att := &CITypeAttribute{
		Name:   "Test",
		Type:   "string",
		Filter: "^[a-zA-Z]+$",
	}

	var err error

	err = format.Validate(att, "ShouldPass")
	if err != nil {
		t.Errorf("Expected string to validate but it did not:\n%s", err.Error())
	}

	err = format.Validate(att, "Should N0T pass!")
	if err == nil {
		t.Errorf("Expected string to fail validation but it passed")
	}

	att.Type = "notstring"
	err = format.Validate(att, "ShouldPass")
	if err == nil {
		t.Errorf("Expected invalid attribute type to fail but it did not")
	}

}

func TestGroupFormat(t *testing.T) {
	format := GetAttributeFormat("group")
	if format == nil {
		t.Errorf("String attribute format does not appear to be registered")
		return
	}

	att := &CITypeAttribute{
		Name:   "Test",
		Type:   "group",
		Filter: "^[a-zA-Z]+$",
	}

	var err error

	err = format.Validate(att, map[string]interface{}{})
	if err != nil {
		t.Errorf("Expected group to validate but it did not:\n%s", err.Error())
	}

	att.Type = "notgroup"
	err = format.Validate(att, map[string]interface{}{})
	if err == nil {
		t.Errorf("Expected invalid attribute type to fail but it did not")
	}
}
