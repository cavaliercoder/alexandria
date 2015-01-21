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
	"fmt"
	"log"
	"os/user"
	"regexp"
	"strings"
)

func GetShortName(name string) string {
	short := strings.ToLower(name)

	// replace all spaces with hyphens
	short = strings.Replace(short, " ", "-", -1)

	// remove all non alphanumerics and non hyphens or underscores
	r := regexp.MustCompile(`[^a-z0-9-_]+`)
	short = r.ReplaceAllString(short, "")

	// Replace multiple hyphens
	r = regexp.MustCompile(`-+`)
	short = r.ReplaceAllString(short, "-")

	return short
}

func IsValidShortName(name string) bool {
	match, err := regexp.MatchString("^[a-z0-9-_]+$", name)
	if err != nil {
		log.Panic(err)
	}
	return match
}

func ExpandPath(path string) string {
	if path[:1] == "~" {
		usr, _ := user.Current()
		path = fmt.Sprintf("%s%s", usr.HomeDir, path[1:])
	}

	return path
}
