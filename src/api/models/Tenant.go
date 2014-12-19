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
package models

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

type Tenant struct {
	model           `bson:",inline"`
	Name  string    `binding:"required"`
	Code  string
}

func (c *Tenant) Init() {
	c.SetCreated()

	shaSum := sha1.Sum(c.Id.([]byte))
	c.Code = strings.ToUpper(fmt.Sprintf("%x-%x-%x", shaSum[0:3], shaSum[3:6], shaSum[7:10]))
}
