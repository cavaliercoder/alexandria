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
	"errors"
	"fmt"
	"regexp"
)

type StringFormat struct{}

func (c *StringFormat) GetName() string {
	return "string"
}

func (c *StringFormat) Validate(att *CITypeAttribute, val interface{}) error {
	if att.Type != c.GetName() {
		return errors.New(fmt.Sprintf("Attribute '%s' is not the correct type", att.Name))
	}

	// Ensure it is a string value
	valStr, ok := val.(string)
	if !ok {
		return errors.New(fmt.Sprintf("Value for '%s' is not a string", att.Name))
	}

	// Validate required attribute
	if att.Required && valStr == "" {
		return errors.New(fmt.Sprintf("Value for '%s' is required", att.Name))
	}

	// Validate length
	if len(valStr) < att.MinLength {
		return errors.New(fmt.Sprintf("Value for '%s' does not meet minimum length requirement of %d characters", att.Name, att.MinLength))
	}

	if att.MaxLength > 0 && len(valStr) > att.MaxLength {
		return errors.New(fmt.Sprintf("Value for '%s' exceeds maximum length requirement of %d characters", att.Name, att.MaxLength))
	}

	// Apply regex filters
	for _, constraint := range att.Filters {
		match, err := regexp.MatchString(constraint, valStr)
		if err != nil {
			return err
		}
		if !match {
			return errors.New(fmt.Sprintf("Value for '%s' does not match the required filter: /%s/", att.Name, constraint))
		}
	}

	return nil
}
