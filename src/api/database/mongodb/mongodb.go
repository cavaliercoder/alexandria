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
package mongodb

import (
	"errors"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"alexandria/api/configuration"
	"alexandria/api/models"
)

var session *MgoSession

type MgoSession struct {
	*mgo.Session
}

func GetConnection() (*MgoSession, error) {
	// Get app configuration
	config, err := configuration.GetConfig()
	if err != nil {
		return nil, err
	}
	dbConfig := &config.Database

	if session == nil {
		// Establish database connection
		dialInfo := mgo.DialInfo{
			Addrs:    dbConfig.Servers,
			Database: dbConfig.Database,
			Timeout:  time.Duration(dbConfig.Timeout * 1000000000),
			Username: dbConfig.Username,
			Password: dbConfig.Password,
		}

		mgosession, err := mgo.DialWithInfo(&dialInfo)
		if err != nil {
			return nil, err
		}

		// Validate connection
		err = mgosession.Ping()
		if err != nil {
			return nil, err
		}

		session = &MgoSession{mgosession}
	}

	return session, nil
}
func (c *MgoSession) IsBootStrapped() (bool, error) {
	// Get app configuration
	config, err := configuration.GetConfig()
	if err != nil {
		return false, err
	}

	// Get db connection
	db, err := GetConnection()
	if err != nil {
		return false, err
	}

	count, err := db.DB(config.Database.Database).C("config").Find(bson.M{}).Count()
	if err != nil {
		return false, err
	}

	if count != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *MgoSession) BootStrap(answers *configuration.Answers) error {
	// Double check we're not bootstrapped
	booted, err := c.IsBootStrapped()
	if err != nil {
		return err
	}
	if booted {
		return errors.New("database is already bootstrapped")
	}

	// Get app configuration
	config, err := configuration.GetConfig()
	if err != nil {
		return err
	}

	// Connect to database
	session, err := GetConnection()
	if err != nil {
		return err
	}

	db := session.DB(config.Database.Database)

	// Create collections and indexes
	db.C("config").Create(&mgo.CollectionInfo{})

	db.C("tenants").Create(&mgo.CollectionInfo{})
	db.C("tenants").EnsureIndex(mgo.Index{Key: []string{"code"}, Unique: true})

	db.C("users").Create(&mgo.CollectionInfo{})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"apiKey"}, Unique: true, Sparse: true})

	// Create default tenant
	tenant := models.Tenant{
		Name: answers.Tenant.Name,
	}
	tenant.Init()
	err = db.C("tenants").Insert(tenant)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created detault tenant '%s' with ID %s", tenant.Name, tenant.Id.Hex())

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
	log.Printf("Created root user '%s %s <%s>' with ID %s", user.FirstName, user.LastName, user.Email, user.Id.Hex())

	// Create configuration entry
	cfgData := models.Config{
		Version:    "1.0.0",
		RootTenant: tenant.Id,
		RootUser:   user.Id,
	}
	cfgData.Init()
	err = db.C("config").Insert(cfgData)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Configuration initialization completed successfully")

	return nil
}
