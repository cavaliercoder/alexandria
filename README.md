# Alexandria CMDB

*A CMDB from the future!!!*

Alexandria CMDB is an open source configuration management database written in [Google Go](https://golang.org/) with a [MongoDB](http://www.mongodb.org/) backend.

This project is in infancy and not ready for deployment. It aims to achieve the following:

* Fast, lightweight and low configuration overhead

* Intuitive and responsive frontend

* Automated data sourcing, transformation and validation

* Vertical and horizontal scalability

* High availability

* Comprehensive RESTful API

* Multitenanted, cloud or on-premise

* Modular and pluggable

* ITIL compliant

## License

Alexandria CMDB - Open source configuration management database
Copyright (C) 2014  Ryan Armstrong (ryan@cavaliercoder.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
    
## Requirements

* [Google Go](https://code.google.com/p/go/) | [Documentation](https://golang.org/doc/)
  
  * [Martini](http://martini.codegangsta.io/) | [Documentation](http://godoc.org/github.com/go-martini/martini)
  
  * [Martini Binding](https://github.com/martini-contrib/binding/)

  * [mgo.v2](https://labix.org/mgo) | [Documentation](http://godoc.org/gopkg.in/mgo.v2)
  
  * [gin](https://github.com/codegangsta/gin) (For development)

## RESTful API Spec

### Generics

For Go `net/http` response codes see: [Package http](http://golang.org/pkg/net/http/#pkg-constants)

* All entities must be retrievable from one or more persistent URIs

* All response bodies must be empty of in JSON format

* On failure, must return on of:

  * `401 Unauthorized` if the end user is not authorized for the request
  
  * `405 Method not allowed` if the request method is not supported for the requested URI

### Create

* Must be a POST request

* The request body should be a JSON formatted entity

```bash
curl -ikX POST -H "Content Type: application/json" -d "{\"key\":\"value\"}" http://<api-url>
```

* On success, must return `201 Created` or `202 Accepted` if the request has been queued

* On success, must return `Location` header with the relative URL of the new entity

* On failure, must return one of:
  
  * `403 Forbidden` if the entity collection is read only
  
  * `406 Not acceptable` if the entity is invalid

  * `409 Conflict` if the entity already exists

### Retrieve

* Must be a `GET` or `HEAD` request

* The request body should be empty

* On success, must return `200` with a JSON formatted response body

* On success, must return an appropriate `cache-control` header

* On failure, must return one of:

  * `404` if the resource is not found

### Update

### Delete