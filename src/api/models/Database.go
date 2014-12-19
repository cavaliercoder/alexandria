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
	"regexp"
	"strings"
)

type Database struct {
	model 				`json:"-" bson:",inline"`
	TenantId  	interface{}	`json:"-"`
	Name	  	string        	`json:"name" binding:"required"`
	ShortName	string		`json:"shortName"`
	Description	string		`json:"description"`
	Backend		string		`json:"-"`
}

func (c *Database) Init() {
	c.SetCreated()
	c.GetShortName()
}

func (c *Database) GetShortName() string {
	if c.ShortName == "" {
		c.ShortName = strings.ToLower(c.Name)
		
		// replace all spaces with hyphens
		c.ShortName = strings.Replace(c.ShortName, " ", "-", -1)
		
		// remove all non alphanumerics and non hyphens
		r := regexp.MustCompile(`[^a-z0-9-]+`)
		c.ShortName = r.ReplaceAllString(c.ShortName, "")
		
		// Replace multiple hyphens
		r = regexp.MustCompile(`-+`)
		c.ShortName = r.ReplaceAllString(c.ShortName, "-")
	}
	
	return c.ShortName
}