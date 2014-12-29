/*
 * Alexandria CMDB - Open source config management database
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
	"net/http"
	"time"
)

type ApiInfo struct {
	Version     string
	InstallDate time.Time   `json:"-"`
	RootUserId  interface{} `json:"-"`
}

func GetApiInfo(res http.ResponseWriter, req *http.Request) {
	var apiInfo ApiInfo
	err := RootDb().C("apiInfo").Find(nil).One(&apiInfo)
	Handle(res, req, err)

	Render(res, req, http.StatusOK, apiInfo)
}
