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

type UserController struct {
    controller
}

func (c *UserController) Init(app *cli.App) error {
    c.app = app
    
    c.app.Commands = append(c.app.Commands, []cli.Command{
        {
            Name: "users",
            Usage: "Create, retrieve, update or delete users",
            Action: c.GetUser,
            Subcommands: []cli.Command{
                {
                    Name: "get",
                    Usage: "get users",
                    Action: c.GetUser,
                },
                {
                    Name: "add",
                    Usage: "add a user",
                    Action: c.AddUser,
                },
                {
                    Name: "update",
                    Usage: "update users",
                },
                {
                    Name: "delete",
                    Usage: "delete users",
                },
            },
        },
    }...)
    
    return nil
}

func (c *UserController) GetUser(context *cli.Context) {
    c.getResource(context, "/users")
}

func (c *UserController) AddUser(context *cli.Context) {
    c.addResource(context, "/users")
}