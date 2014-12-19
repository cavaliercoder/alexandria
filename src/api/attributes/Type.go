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
 * package controllers
 */

package attributes

type Type interface {
	GetName() string
	GetDescription() string
	GetValue() string
	SetValue(string) error
}

type typ struct {
	value string
}

func (c *typ) GetValue() string {
	return c.value
}

func (c *typ) SetValue(value string) error {
	c.value = value
	return nil
}

type TypeMap map[string]Type

var types TypeMap

func GetAttributeTypes() TypeMap {
	if len(types) == 0 {
		typs := []Type{&String{}}

		for _, typ := range typs {
			types[typ.GetName()] = typ
		}
	}

	return types
}
