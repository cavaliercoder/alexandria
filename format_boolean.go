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
	"strings"
)

type BooleanFormat struct{}

func (c *BooleanFormat) GetName() string {
	return "boolean"
}

func (c *BooleanFormat) Validate(att *CITypeAttribute, val *interface{}) error {
	if att.Type != c.GetName() {
		return errors.New(fmt.Sprintf("Attribute '%s' is not the correct type", att.Name))
	}

	// Parse a string
	if str, ok := (*val).(string); ok {
		str = strings.ToLower(str)
		truthies := []string{
			"true",
			"yes",
			"y",
			"1",
		}

		for _, t := range truthies {
			if str == t {
				*val = true
				return nil
			}
		}

		falsies := []string{
			"false",
			"no",
			"n",
			"0",
		}

		for _, f := range falsies {
			if str == f {
				*val = false
				return nil
			}
		}

		return errors.New(fmt.Sprintf("Value '%s' for '%s' is not a valid boolean value", str, att.Name))
	}

	// Parse a number
	if num, ok := (*val).(int); ok {
		if num > 0 {
			*val = true
		} else {
			*val = false
		}

		return nil
	}

	// Parse as native boolean
	_, ok := (*val).(bool)
	if !ok {
		return errors.New(fmt.Sprintf("Value for '%s' is not a valid boolean", att.Name))
	}

	return nil
}
