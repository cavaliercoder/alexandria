# Alexandria

*A CMDB from the future!!!*

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