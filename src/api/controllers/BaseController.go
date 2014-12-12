package controllers

import (
    "github.com/go-martini/martini"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "encoding/json"
    "net/http"
    "log"
    "alexandria/api/models"
)

type baseController struct {
    m *martini.ClassicMartini
    db *mgo.Database
}

func (c baseController) RenderJson(w http.ResponseWriter, v interface{}) {
    if v == nil {
        var empty []struct{}
        v = empty
    }
    
    json, err := json.Marshal(v)
    c.Handle(err)
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(json)
}

func (c baseController) Handle(err error) {
    if err != nil {
        log.Panic(err)
    }
}

func (c baseController) GetEntities(collection string, w http.ResponseWriter) {    
    dbcollection := c.db.C(collection)
    var results []interface{}
    err := dbcollection.Find(bson.M{}).All(&results)
    c.Handle(err)
    
    c.RenderJson(w, results)
}

func (c baseController) AddEntity(collection string, uri string, model models.BaseModel, w http.ResponseWriter) {
    
    // Insert new user
    model.SetCreated()
    
    err := c.db.C(collection).Insert(&model)
    c.Handle(err)    
    w.Header().Set("Location", uri)
    
    c.RenderJson(w, model)
}