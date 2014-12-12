package controllers

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/binding"
    "gopkg.in/mgo.v2"
    "net/http"
    "alexandria/api/models"
    "fmt"
)

type tenantController struct {
    baseController
}

func NewTenantController(m *martini.ClassicMartini, db *mgo.Database) (*tenantController, error) {
    c := new(tenantController)    
    c.m = m
    c.db = db    
    
    // Add routes
    m.Get("/tenants", c.GetTenants)
    //m.Get("/tenants/:email", c.GetUserByEmail)
    m.Post("/tenants", binding.Bind(models.Tenant{}), c.AddTenant)
    
    return c, nil
}

func (c tenantController) GetTenants(w http.ResponseWriter) {    
    c.GetEntities("tenants", w)
}

func (c tenantController) AddTenant(tenant models.Tenant, w http.ResponseWriter) {
    c.AddEntity("tenants", fmt.Sprintf("/tenants/%s", tenant.Id), &tenant, w)
}