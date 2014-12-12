package controllers

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/binding"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "net/http"
    "alexandria/api/models"
    "log"
    "fmt"
)

type userController struct {
    baseController
}

func NewUserController(m *martini.ClassicMartini, db *mgo.Database) (*userController, error) {
    c := new(userController)    
    c.m = m
    c.db = db    
    
    // Add routes
    m.Get("/users", c.GetUsers)
    m.Get("/users/:email", c.GetUserByEmail)
    m.Post("/users", binding.Bind(models.User{}), c.AddUser)
  
    // Initialize database
    c.db.C("users").Create(&mgo.CollectionInfo{})
    c.db.C("users").EnsureIndex(mgo.Index{ Key: []string{"Email"}, Unique: true})
    
    return c, nil
}

func (c userController) GetUsers(w http.ResponseWriter) {    
    c.GetEntities("users", w)
}

func (c userController) GetUserByEmail(params martini.Params, w http.ResponseWriter) {
    var user models.User
    err := c.db.C("users").Find(bson.M{"email": params["email"]}).One(&user)
    c.Handle(err)
    
    c.RenderJson(w, user)
}

func (c userController) AddUser(user models.User, w http.ResponseWriter) {
    // Make sure user doesn't already exist
    count, err := c.db.C("users").Find(bson.M{"email": user.Email}).Count()
    c.Handle(err)    
    if count > 0 {
        w.WriteHeader(http.StatusConflict)
        log.Panic(fmt.Sprintf("A user account already exists for email %s", user.Email))
    }
    
    c.AddEntity("users", fmt.Sprintf("/users/%s", user.Email), &user, w)
}