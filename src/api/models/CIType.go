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
package models

import ()

type CIType struct {
	tenantedModel `bson:",inline"`

	InheritFrom    interface{}	`json:"inheritFrom"`
	Name        string `json:"name" binding:"required"`
	ShortName   string `json:"shortName"`
	Description string `json:"description"`
	Attributes  []CIAttribute `json:"attributes"`
}

func (c *CIType) Init() {
	c.SetCreated()
}
