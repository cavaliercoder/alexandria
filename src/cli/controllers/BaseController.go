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
  
  "bufio"
  "net/http"
  "fmt"
  "io"
  "strings"
  
  . "alexandria/cli/application"
)

type BaseController interface {
    Init(app *cli.App) error
}

type baseController struct {
    app     *cli.App
}

func (c baseController) ApiRequest(context *cli.Context, method string, path string, body io.Reader) (*http.Request, *http.Response, error) {
    url := fmt.Sprintf("%s%s", context.GlobalString("url"), path)
    
    DPrint(fmt.Sprintf("API Request: %s %s", method, url))
    
    client := new(http.Client)    
    req, err := http.NewRequest(strings.ToUpper(method), url, body)
    if err != nil { return req, &http.Response{}, err }
    defer req.Body.Close()
    
    req.Header.Add("Content-type", "application/json")
    req.Header.Add("X-Auth-Token", context.String("api-key"))
    req.Header.Add("User-Agent", "Alexandria CMDB CLI")
        
    var res *http.Response
    res, err = client.Do(req)
    
    return req, res, err
}

func (c baseController) DumpHttpError(req *http.Request, res *http.Response) {    
    // Print request body
    if req != nil {
        DPrint("<><><> Request <><><>")
        
        for k := range req.Header {
            DPrint(fmt.Sprintf("Header %s: %s", k, req.Header[k][0]))
        }
        
        DPrint("Body:")
        buf := bufio.NewReader(req.Body)        
        var line []byte
        var err error = nil
        for err == nil {
           line, _, err = buf.ReadLine()           
           DPrint(line)
        }
    }
    
    if res != nil {
        DPrint("<><><> Response <><><>")
        DPrint(fmt.Sprintf("Status: HTTP %s", res.Status))
        
        for k := range res.Header {
            DPrint(fmt.Sprintf("Header %s: %s", k, res.Header[k][0]))
        }
        
        DPrint("Body:")
        
        // Print request body
        buf := bufio.NewReader(res.Body)        
        var line []byte
        var err error = nil
        for err == nil {
           line, _, err = buf.ReadLine()           
           DPrint(line)
        }
    }
}