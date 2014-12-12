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
package main

import (
  "os"
  "github.com/codegangsta/cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "alex"
    app.Usage = "Alexandria CMDB CLI"
    app.Version = "1.0.0"
  
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "u, url",
            Value: "http://localhost:3000",
            Usage: "Specify the API base URL",
            EnvVar: "ALEX_API_URL",
        },
        cli.BoolFlag{
            Name:  "verbose",
            Usage: "Show more output",
        },
    }
    
    // Commands
    app.Commands = []cli.Command{
    {
        Name: "build",
        Flags: []cli.Flag{
            cli.BoolFlag{
                Name:  "no-cache",
                Usage: "Do not use cache when building the image.",
            },
        },
        Usage:  "Build or rebuild services",
        Action: CmdBuild,
    }}

    app.Action = func(c *cli.Context) {
        println("Hello friend!")
    }

    app.Run(os.Args)
}

func CmdBuild(c *cli.Context) {
    println("Yep!")
}