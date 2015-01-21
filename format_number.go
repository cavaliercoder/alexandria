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
	"errors"
	"fmt"
	"strconv"
)

type NumberFormat struct{}

func (c *NumberFormat) GetName() string {
	return "number"
}

func (c *NumberFormat) Validate(att *CITypeAttribute, val *interface{}) error {
	if att.Type != c.GetName() {
		return errors.New(fmt.Sprintf("Attribute '%s' is not the correct type", att.Name))
	}

	// If the user submitted the value as a string, it must be converted
	if str, ok := (*val).(string); ok {
		f64, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return errors.New(fmt.Sprintf("Value '%v' for attribute '%s' is not a valid number", *val, att.Name))
		}

		// Update the user submitted value
		*val = f64
	}

	// Number is always stored as a 64bit floating point integer (float64)
	if f64, ok := (*val).(float64); ok {
		if att.MinValue != 0 && f64 < att.MinValue {
			return errors.New(fmt.Sprintf("Value '%v' for attribute '%s' does not exceed minimum value '%v'", *val, att.Name, att.MinValue))
		}

		if att.MaxValue != 0 && f64 > att.MaxValue {
			return errors.New(fmt.Sprintf("Value '%v' for attribute '%s' exceeds maximum value '%v'", *val, att.Name, att.MinValue))
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Value '%v' for attribute '%s' is not a valid number", *val, att.Name))
}
