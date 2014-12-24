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
package main

import (
	"regexp"
	"strings"
)

func GetShortName(name string) string {
	short := strings.ToLower(name)

	// replace all spaces with hyphens
	short = strings.Replace(short, " ", "-", -1)

	// remove all non alphanumerics and non hyphens
	r := regexp.MustCompile(`[^a-z0-9-]+`)
	short = r.ReplaceAllString(short, "")

	// Replace multiple hyphens
	r = regexp.MustCompile(`-+`)
	short = r.ReplaceAllString(short, "-")

	return short
}
