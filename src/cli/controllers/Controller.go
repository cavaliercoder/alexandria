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
    
    // Validate URL and API Key
    if url == "" { Die("API base URL not specified\n") }
    if apiKey == "" { Die("API authentication key not specified\n") }
    
    // Formulate request URL
    url = fmt.Sprintf("%s%s?pretty=true", url, path)    
    DPrint(fmt.Sprintf("API Request: %s %s", method, url))
    
    // Create a HTTP client that does not follow redirects
    // This allows 'Location' headers to be printed to the CLI
    client := http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return errors.New("Never follow redirects")
        },
    }
    req, err := http.NewRequest(strings.ToUpper(method), url, body)
    if err != nil { return nil, err }
    
    // Add request headers
    req.Header.Add("Content-type", "application/json")
    req.Header.Add("X-Auth-Token", apiKey)
    req.Header.Add("User-Agent", "Alexandria CMDB CLI")
    
    // Submit the request  
    res, err := client.Do(req)
    
    return res, err
}

func (c *controller) ApiResult(res *http.Response) {
    defer res.Body.Close()
    io.Copy(os.Stdout, res.Body)
    fmt.Println()
}

func (c *controller) ApiError(res *http.Response) {
        fmt.Fprintf(os.Stderr, "%s\n", res.Status)
        //io.Copy(os.Stderr, res.Body)
        //fmt.Fprintf(os.Stderr, "\n")
        os.Exit(1)
}

func (c *controller) getResource(context *cli.Context, path string) {
    // Get requested resource ID from first command argument
    id := context.Args().First()
    
    var err error
    var res *http.Response
    
    if id != "" {
        // Get one by id
        path = fmt.Sprintf("%s/%s", path, id)
    }
    
    res, err = c.ApiRequest(context, "GET", path, nil)
    if err != nil { Die(err) }
    
    switch res.StatusCode {
        case http.StatusOK:
            c.ApiResult(res)
        case http.StatusNotFound:
            Die(fmt.Sprintf("No such resource found at %s", path))
        default:
            c.ApiError(res)
    }
}

func (c *controller) addResource(context *cli.Context, path string, resource string) {
    // Decode the resource from STDIN or from the first command argument?
    var input io.Reader
    if context.GlobalBool("stdin") {
        input = os.Stdin
    } else {
        if resource == "" {
            resource = context.Args().First()
        }
        
        input = strings.NewReader(resource)
    }
    
    res, err := c.ApiRequest(context, "POST", path, input)
    defer res.Body.Close()
    if err != nil { Die(err) }
    
    if res.StatusCode == http.StatusCreated {
        fmt.Printf("Created %s\n", res.Header.Get("Location"))
    } else {
        c.ApiError(res)
    }
}

func (c *controller) deleteResource(context *cli.Context, path string) {
     // Get requested resource ID from first command argument
    id := context.Args().First()
    
    var err error
    var res *http.Response
    
    if id == "" {
        Die("No user specified")
    }
    
    path = fmt.Sprintf("%s/%s", path, id)
    
    res, err = c.ApiRequest(context, "DELETE", path, nil)
    if err != nil { Die(err) }
    
    switch res.StatusCode {
        case http.StatusNoContent:
            fmt.Printf("Deleted %s\n", path)
        case http.StatusNotFound:
            Die(fmt.Sprintf("No such resource found at %s", path))
        default:
            c.ApiError(res)
    }   
}