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
	"testing"
)

func TestFormatFactory(t *testing.T) {
	if format := GetAttributeFormat("i_dont_exist"); format != nil {
		t.Errorf("Expected nil when requesting nonexistant Attribute Format but got: %#v", format)
	}
}

func TestStringFormat(t *testing.T) {
	var err error

	format := GetAttributeFormat("string")
	if format == nil {
		t.Errorf("String attribute format does not appear to be registered")
		return
	}

	att := &CITypeAttribute{
		Name: "Test",
		Type: "notstring",
	}

	var val interface{}

	// Test invalid attribute type
	val = "ShouldPass"
	err = format.Validate(att, &val)
	if err == nil {
		t.Errorf("Expected invalid attribute type to fail but it did not")
	}
	att.Type = "string"

	// Test filters
	att.Filters = []string{"^[a-zA-Z]+$", "^ShouldPass$"}
	err = format.Validate(att, &val)
	if err != nil {
		t.Errorf("Expected string to validate but it did not:\n%s", err.Error())
	}

	val = "ShouldNotPass"
	err = format.Validate(att, &val)
	if err == nil {
		t.Errorf("Expected string to fail validation but it passed")
	}
	att.Filters = nil

	// Test required value
	att.Required = true
	val = ""
	err = format.Validate(att, &val)
	if err == nil {
		t.Errorf("Expected string to fail with a required value but it passed")
	}
	att.Required = false

	// Test minimum length
	att.MinLength = 10
	val = "too short"
	err = format.Validate(att, &val)
	if err == nil {
		t.Errorf("Expected string to fail minimum length requirement but it passed")
	}
	att.MinLength = 0

	// Test maximum length
	att.MaxLength = 7
	val = "too long"
	err = format.Validate(att, &val)
	if err == nil {
		t.Errorf("Expected string to fail maximum length requirement but it passed")
	}
	att.MaxLength = 0

	// Test multiple
	att.MinLength = 17
	att.MaxLength = 17
	att.Required = true
	att.Filters = []string{"^Lorem Ipsum D0lor$", "^[LoremIpsuD0l ]+$", "^[a-zA-Z0-9 ]+$"}
	val = "Lorem Ipsum D0lor"
	err = format.Validate(att, &val)
	if err != nil {
		t.Errorf("Expected string to pass multiple requirements but it failed with:\n    %s", err.Error())
	}
}

func TestGroupFormat(t *testing.T) {
	format := GetAttributeFormat("group")
	if format == nil {
		t.Errorf("String attribute format does not appear to be registered")
		return
	}

	att := &CITypeAttribute{
		Name:    "Test",
		Type:    "group",
		Filters: []string{"^[a-zA-Z]+$"},
	}

	var err error
	var val interface{}

	val = map[string]interface{}{}
	err = format.Validate(att, &val)
	if err != nil {
		t.Errorf("Expected group to validate but it did not:\n%s", err.Error())
	}

	att.Type = "notgroup"
	err = format.Validate(att, &val)
	if err == nil {
		t.Errorf("Expected invalid attribute type to fail but it did not")
	}
}

func TestBooleanFormat(t *testing.T) {
	format := GetAttributeFormat("boolean")
	if format == nil {
		t.Errorf("Boolean attribute format does not appear to be registered")
		return
	}

	var err error
	att := &CITypeAttribute{
		Name: "Test",
		Type: "boolean",
	}

	var val interface{}

	// Test native booleans
	val = true
	err = format.Validate(att, &val)
	if err != nil {
		t.Errorf("Expected boolean attribute to validate but it did not")
	}

	val = false
	err = format.Validate(att, &val)
	if err != nil {
		t.Errorf("Expected boolean attribute to validate but it did not")
	}

	// Test strings
	val = "TRUE"
	err = format.Validate(att, &val)
	if err != nil {
		t.Error("Expected boolean string to validate but it did not")
	}
	if val != true {
		t.Errorf("Expected native true; got: %#v", val)
	}

	val = "False"
	err = format.Validate(att, &val)
	if err != nil {
		t.Error("Expected boolean string to validate but it did not")
	}
	if val != false {
		t.Errorf("Expected native true; got: %#v", val)
	}

	val = "Not a truthy or falsy"
	err = format.Validate(att, &val)
	if err == nil {
		t.Error("Expected boolean string to validate but it did not")
	}

	// Test numbers
	val = 10
	err = format.Validate(att, &val)
	if err != nil {
		t.Error("Expected boolean number to validate but it did not")
	}
	if val != true {
		t.Errorf("Expected native true; got: %#v", val)
	}

	val = -10
	err = format.Validate(att, &val)
	if err != nil {
		t.Error("Expected boolean number to validate but it did not")
	}
	if val != false {
		t.Errorf("Expected native true; got: %#v", val)
	}

}
