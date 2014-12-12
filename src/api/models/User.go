package models

import (
    "gopkg.in/mgo.v2/bson"
)

type User struct {
    baseModel                       `bson:",inline"`
    
    TenantId        bson.ObjectId   `json:"-"`
    FirstName       string          `json:"firstName"`
    LastName        string          `json:"lastName"`
    Email           string          `json:"email" binding:"required"`
}