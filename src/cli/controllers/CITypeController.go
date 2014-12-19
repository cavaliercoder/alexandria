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
package controllers

import (
	"github.com/codegangsta/cli"
)

type CITypeController struct {
	controller
}

func (c *CITypeController) Init(app *cli.App) error {
	c.app = app

	c.app.Commands = append(c.app.Commands, []cli.Command{
		{
			Name:   "types",
			Usage:  "Create, retrieve, update or delete CI types",
			Action: c.GetCIType,
			Subcommands: []cli.Command{
				{
					Name:   "get",
					Usage:  "get types",
					Action: c.GetCIType,
				},
				{
					Name:   "add",
					Usage:  "add types",
					Action: c.AddCIType,
				},
				{
					Name:  "update",
					Usage: "update types",
				},
				{
					Name:  "delete",
					Usage: "delete types",
				},
			},
		},
	}...)

	return nil
}

func (c *CITypeController) GetCIType(context *cli.Context) {
	c.getResource(context, "/types")
}

func (c *CITypeController) AddCIType(context *cli.Context) {
	c.addResource(context, "/types", "")
}
