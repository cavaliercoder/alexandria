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
  "errors"
  "fmt"
  "io"
  "net/http"
  "os"
  "strings"
  
  . "alexandria/cli/application"
)

type Controller interface {
    Init(app *cli.App) error
}

type controller struct {
    app     *cli.App
}

func (c *controller) ApiRequest(context *cli.Context, method string, path string, body io.Reader) (*http.Response, error) {
    url := context.GlobalString("url")
    apiKey := context.GlobalString("api-key")
    
    if url == "" {
        fmt.Fprintf(os.Stderr, "API base URL not specified\n")
        os.Exit(1)
    }
    
    if apiKey == "" {
        fmt.Fprintf(os.Stderr, "API authentication key not specified\n")
        os.Exit(1)
    }
    
    url = fmt.Sprintf("%s%s?pretty=true", url, path)    
    DPrint(fmt.Sprintf("API Request: %s %s", method, url))
    
    client := http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return errors.New("Never follow redirects")
        },
    }
    req, err := http.NewRequest(strings.ToUpper(method), url, body)
    if err != nil { return nil, err }
    
    req.Header.Add("Content-type", "application/json")
    req.Header.Add("X-Auth-Token", apiKey)
    req.Header.Add("User-Agent", "Alexandria CMDB CLI")
        
    var res *http.Response
    res, err = client.Do(req)
    
    return res, err
}

func (c *controller) ApiResult(res *http.Response) {
    defer res.Body.Close()
    io.Copy(os.Stdout, res.Body)
    fmt.Println()
}

func (c *controller) ApiError(res *http.Response) {
        fmt.Fprintf(os.Stderr, "%s\n", res.Status)
        io.Copy(os.Stderr, res.Body)
        fmt.Fprintf(os.Stderr, "\n")
        os.Exit(1)
}

func (c *controller) DumpHttpError(req *http.Request, res *http.Response) {    
    // Print request body
    if req != nil {
        defer req.Body.Close()
        
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