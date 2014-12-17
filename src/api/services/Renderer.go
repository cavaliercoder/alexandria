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
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
)

type Renderer struct {
	*http.Request
	http.ResponseWriter
}

// Wire the service
func RendererService() martini.Handler {
	return func(req *http.Request, res http.ResponseWriter, c martini.Context) {
		r := &Renderer{req, res}

		c.Map(r)
	}
}

func (c *Renderer) Handle(err error) {
	if err == nil {
		return
	}

	switch err.Error() {
	case "not found":
		c.WriteHeader(http.StatusNotFound)
	default:
		log.Panic(err)
	}
}

func (c *Renderer) NotFound() {
	c.WriteHeader(http.StatusNotFound)
}

func (c *Renderer) Render(status int, v interface{}) {
	format := c.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	switch format {
	case "json":
		c.JSON(status, v)

	default:
		log.Panic(fmt.Sprintf("Unsupported output format: %s", format))
	}
}

func (c *Renderer) JSON(status int, v interface{}) {
	if v == nil {
		v = new(struct{})
	}

	var err error
	var data []byte
	if c.URL.Query().Get("pretty") == "true" {
		data, err = json.MarshalIndent(v, "", "    ")
	} else {
		data, err = json.Marshal(v)
	}
	if err != nil {
		log.Panic(err)
	}

	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.WriteHeader(status)

	c.ResponseWriter.Write(data)
}
