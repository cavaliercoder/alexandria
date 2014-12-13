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
  
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
)

type ConfigController struct {
    baseController    
}

func (c ConfigController) Init(app *cli.App) error {
    c.app = app
    
    c.app.Commands = append(c.app.Commands, []cli.Command{
        {
            Name: "config",
            Usage: "manage server configuration",
            Subcommands: []cli.Command{
                {
                    Name: "get",
                    Usage: "get configuration",
                },
                {
                    Name: "init",
                    Usage: "initialize a new server",
                    Action: c.InitConfig,
                },
            },
        },
    }...)
    
    return nil
}

func (c ConfigController) InitConfig(context *cli.Context) {    
    _, res, err := c.ApiRequest(context, "GET", "/config?init=true", nil)
    if err != nil { log.Panic(err) }
    defer res.Body.Close()
    
    switch res.StatusCode {
        case http.StatusOK:
            io.Copy(os.Stdout, res.Body)
            
        case http.StatusNotFound:
            fmt.Fprintf(os.Stderr,"Server configuration is already intialized\n")
            os.Exit(1)
            
        default:
            fmt.Fprintf(os.Stderr, "%s\n", res.Status)
            io.Copy(os.Stderr, res.Body)
            os.Exit(1)
    }
}