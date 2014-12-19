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
  "log"
  "net/http"
  "os"
)

type ConfigController struct {
    controller    
}

func (c *ConfigController) Init(app *cli.App) error {
    c.app = app
    
    c.app.Commands = append(c.app.Commands, []cli.Command{
        {
            Name: "config",
            Usage: "manage server configuration",
            Subcommands: []cli.Command{
                {
                    Name: "get",
                    Usage: "get configuration",
                    Action: c.GetConfig,
                },
                {
                    Name: "init",
                    Usage: "initialize a new server",
                    Action: c.InitConfig,
                },
                {
                    Name: "destroy",
                    Usage: "reset factory defaults",
                    Action: c.ClearConfig,
                },
            },
        },
    }...)
    
    return nil
}

func (c *ConfigController) GetConfig(context *cli.Context) {
    res, err := c.ApiRequest(context, "GET", "/config", nil)
    if err != nil { log.Panic(err) }
    
    c.ApiResult(res)
}

func (c *ConfigController) InitConfig(context *cli.Context) {    
    res, err := c.ApiRequest(context, "POST", "/config/actions/initialize", nil)
    if err != nil { log.Panic(err) }
    
    switch res.StatusCode {
        case http.StatusCreated, http.StatusOK :
            c.ApiResult(res)
            
        case http.StatusNotFound:
            fmt.Fprintf(os.Stdout, "Configuration is already initialized.\n")
            os.Exit(1)
            
        default:
            c.ApiError(res)
    }
}

func (c *ConfigController) ClearConfig(context *cli.Context) {
    res, err := c.ApiRequest(context, "POST", "/config/actions/destroy", nil)
    if err != nil { log.Panic(err) }
    defer res.Body.Close()
    
    if(res.StatusCode == http.StatusOK) {
        fmt.Fprintf(os.Stdout, "Configuration destroyed.\n")
    } else {
        c.ApiError(res)
    }
}