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

type AttributeFormat interface {
	GetName() string
	Validate(*CITypeAttribute, *interface{}) error
}

var formatMap map[string]AttributeFormat

func GetAttributeFormat(name string) AttributeFormat {
	if formatMap == nil {
		// Initialize the map
		formatMap = map[string]AttributeFormat{}
		formats := []AttributeFormat{
			&GroupFormat{},
			&StringFormat{},
			&NumberFormat{},
			&BooleanFormat{},
			&TimeStampFormat{},
		}

		for _, format := range formats {
			formatMap[format.GetName()] = format
		}
	}

	return formatMap[name]
}
