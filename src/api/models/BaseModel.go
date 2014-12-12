package models

import (
    "gopkg.in/mgo.v2/bson"
    "time"
)

type BaseModel interface {
    SetCreated()
    SetModified()
}

type baseModel struct {
    Id          bson.ObjectId      `json:"-" bson:"_id,omitempty"`
    Created     time.Time          `json:"-" bson:"created"`
    Modified    time.Time          `json:"-" bson:"modified"`
}

func (c *baseModel) SetCreated()  {
    
    if c.Id.Hex() == "" {
        c.Id = bson.NewObjectId()
    }
    
    if c.Created.IsZero() {
        now := time.Now()
        c.Created = now
        c.Modified = now
    }
}

func (c *baseModel) SetModified() {
    c.Modified = time.Now()
}