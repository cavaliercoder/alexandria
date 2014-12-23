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

package attributes

import (
	"errors"
	"fmt"
	"strings"
)

type Boolean struct {
	typ
}

func (c *Boolean) GetName() string {
	return "Boolean"
}

func (c *Boolean) GetDescription() string {
	return "True/False binary value"
}

func (c *Boolean) SetValue(value string) error {
	switch strings.ToLower(value) {
	case "true", "1", "yes", "y":
		c.value = "True"
	case "false", "0", "no", "n", "null", "nil":
		c.value = "False"
	default:
		return errors.New(fmt.Sprintf("Invalid boolean value: %s", value))
	}

	return nil
}
