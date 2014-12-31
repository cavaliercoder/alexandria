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
)

type GroupFormat struct{}

func (c *GroupFormat) GetName() string {
	return "group"
}

func (c *GroupFormat) Validate(att *CITypeAttribute, val interface{}) error {
	if att.Type != c.GetName() {
		return errors.New(fmt.Sprintf("Attribute '%s' is not the correct type", att.Name))
	}

	_, ok := val.(map[string]interface{})
	if !ok {
		return errors.New(fmt.Sprintf("Value for '%s' is not an attribute group", att.Name))
	}

	return nil
}
