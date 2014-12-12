package models

type Tenant struct {
    baseModel               `bson:",inline"`
    Name        string      `json:"name"`
}