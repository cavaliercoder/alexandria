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
  "strings"
)

type TenantController struct {
    baseController    
}

func (c TenantController) Init(app *cli.App) error {
    c.app = app
    
    c.app.Commands = append(c.app.Commands, []cli.Command{
        {
            Name: "tenants",
            Usage: "Create, retrieve, update or delete tenants",
            Action: c.GetTenant,
            Subcommands: []cli.Command{
                {
                    Name: "get",
                    Usage: "get tenants",
                    Action: c.GetTenant,
                },
                {
                    Name: "add",
                    Usage: "add tenants",
                    Action: c.AddTenant,
                },
                {
                    Name: "update",
                    Usage: "update tenants",
                },
                {
                    Name: "delete",
                    Usage: "delete tenants",
                },
            },
        },
    }...)
    
    return nil
}

func (c TenantController) GetTenant(context *cli.Context) {
    id := context.Args().First()
    
    var err error
    var res *http.Response
    
    if id == "" {
        _, res, err = c.ApiRequest(context, "GET", "/tenants", nil)
    } else {
        _, res, err = c.ApiRequest(context, "GET", fmt.Sprintf("/tenants/%s", id), nil)
    }
    
    if err != nil {
        log.Panic(err)
    }
    
    c.ApiResult(res)
}

func (c TenantController) AddTenant(context *cli.Context) {
    var input io.Reader
    if context.GlobalBool("stdin") {
        input = os.Stdin
    } else {
        input = strings.NewReader(context.Args().First())
    }
    
    _, res, err := c.ApiRequest(context, "POST", "/tenants", input)
    if err != nil { log.Panic(err) }
    defer res.Body.Close()
    
    if res.StatusCode == http.StatusCreated {
        fmt.Printf("Created %s\n", res.Header.Get("Location"))
    } else {
        c.ApiError(res)
    }
}