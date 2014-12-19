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
package database

import (
	"errors"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"alexandria/api/configuration"
	"alexandria/api/models"
)

type MongoDriver struct {
	session	*mgo.Session
	rootDB	*mgo.Database
	config 	*configuration.DatabaseConfig
}

func (c *MongoDriver) Connect(config *configuration.DatabaseConfig) error {
	if c.session == nil || c.config != config {
		if c.session != nil { c.session.Close() }
		
		var err error
		
		// Establish database connection
		dialInfo := mgo.DialInfo{
			Addrs:    config.Servers,
			Database: config.Database,
			Timeout:  time.Duration(config.Timeout * 1000000000),
			Username: config.Username,
			Password: config.Password,
		}

		log.Printf("MongoDB: Connecting to %s (%s)...", config.Servers, config.Database)
		c.session, err = mgo.DialWithInfo(&dialInfo)
		if err != nil {
			return err
		}
		
		// enable error checking
		c.session.SetSafe(&mgo.Safe{})

		// Validate connection
		log.Printf("MongoDB: Validating connection...")
		err = c.session.Ping()
		if err != nil {
			return err
		}
		
		c.config = config
		c.rootDB = c.session.DB(config.Database)
	}

	return nil
}

func (c *MongoDriver) Clone () (Driver, error) {
	clone := new(MongoDriver)
	clone.session = c.session.Clone()
	clone.config = c.config
	clone.rootDB = clone.session.DB(c.config.Database)
	
	return clone, nil
}

func (c *MongoDriver) Close() error {
	c.session.Close()
	return nil
}

func (c *MongoDriver) IsBootStrapped() (bool, error) {
	count, err := c.rootDB.C("config").Find(nil).Count()
	if err != nil {
		return false, err
	}

	if count != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *MongoDriver) BootStrap(answers *configuration.Answers) error {
	// Double check we're not bootstrapped
	booted, err := c.IsBootStrapped()
	if err != nil {
		return err
	}
	if booted {
		return errors.New("database is already bootstrapped")
	}
	
	db:= c.rootDB
	
	// Create collections and indexes
	log.Printf("Creating collections and indexes...")
	db.C("config").Create(&mgo.CollectionInfo{})

	db.C("tenants").Create(&mgo.CollectionInfo{})
	db.C("tenants").EnsureIndex(mgo.Index{Key: []string{"code"}, Unique: true})

	db.C("users").Create(&mgo.CollectionInfo{})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"apikey"}, Unique: true, Sparse: true})

	db.C("databases").Create(&mgo.CollectionInfo{})
	db.C("databases").EnsureIndex(mgo.Index{Key: []string{"tenantid"}, Unique: false})
	
	// Create default tenant
	tenant := models.Tenant{
		Name: answers.Tenant.Name,
	}
	tenant.Init()
	err = db.C("tenants").Insert(tenant)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created detault tenant '%s' with ID %s", tenant.Name, tenant.Id.(bson.ObjectId).Hex())

	// Create root user
	user := models.User{
		FirstName: answers.User.FirstName,
		LastName:  answers.User.LastName,
		Email:     answers.User.Email,
	}
	user.Init()
	user.TenantId = tenant.Id

	err = db.C("users").Insert(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created root user '%s %s <%s>' with ID %s", user.FirstName, user.LastName, user.Email, user.Id.(bson.ObjectId).Hex())

	// Create configuration entry
	cfgData := models.Config{
		Version:    "1.0.0",
		RootTenant: tenant.Id.(bson.ObjectId),
		RootUser:   user.Id.(bson.ObjectId),
	}
	cfgData.Init()
	err = db.C("config").Insert(cfgData)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Configuration initialization completed successfully")

	return nil
}

func (c *MongoDriver) IdToString(id interface{}) string {
	return id.(bson.ObjectId).Hex()	
}

func (c *MongoDriver) CreateDatabase(database string) error {
	return nil
}

func (c *MongoDriver) DeleteDatabase(database string) error {
	err := c.session.DB(database).DropDatabase()
	
	return err
}

func (c *MongoDriver) GetAll(collection string, filter M, results interface{}) error {
	err := c.rootDB.C(collection).Find(filter).All(results)
	return err
}

func (c *MongoDriver) GetOne(collection string, filter M, result interface{}) error {
	err := c.rootDB.C(collection).Find(filter).One(result)
	return err
}

func (c *MongoDriver) GetOneById(collection string, id interface{}, result interface{}) error {
	err := c.rootDB.C(collection).FindId(id).One(result)
	return err
}

func (c *MongoDriver) Insert(collection string, item models.Model) error {
	if item.GetId() == nil || c.IdToString(item.GetId()) == "" {
		item.SetId(bson.NewObjectId())
	}
	
	err := c.rootDB.C(collection).Insert(item)
	return err
}

func (c *MongoDriver) Delete(collection string, filter M) error {
	err := c.rootDB.C(collection).Remove(filter)
	return err
}