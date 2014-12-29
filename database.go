/*
 * Alexandria CMDB - Open source common.management database
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
package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"reflect"
	"time"
)

type M map[string]interface{}

var dbSession *mgo.Session

func DbConnect() *mgo.Session {
	if dbSession == nil {
		config, err := GetConfig()
		if err != nil {
			log.Panic(err)
		}

		// Establish database connection
		dialInfo := mgo.DialInfo{
			Addrs:    config.Database.Servers,
			Database: config.Database.Database,
			Timeout:  time.Duration(config.Database.Timeout * 1000000000),
			Username: config.Database.Username,
			Password: config.Database.Password,
		}

		log.Printf("MongoDB: Connecting to %s (%s)...", config.Database.Servers, config.Database.Database)
		dbSession, err = mgo.DialWithInfo(&dialInfo)
		if err != nil {
			log.Panic(err)
		}

		// enable error checking
		dbSession.SetSafe(&mgo.Safe{})

		// Validate connection
		log.Printf("MongoDB: Validating connection...")
		err = dbSession.Ping()
		if err != nil {
			log.Panic(err)
		}
	}

	return dbSession.Clone()
}

func Db(name string) *mgo.Database {
	session := DbConnect()
	return session.DB(name)
}

func RootDb() *mgo.Database {
	config, err := GetConfig()
	if err != nil {
		log.Panic(err)
	}

	session := DbConnect()
	return session.DB(config.Database.Database)
}

func IsBootStrapped() (bool, error) {
	count, err := RootDb().C("config").Find(nil).Count()
	if err != nil {
		return false, err
	}

	if count != 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func BootStrap(answers *Answers) error {
	// Double check we're not bootstrapped
	booted, err := IsBootStrapped()
	if err != nil {
		return err
	}
	if booted {
		return errors.New("database is already bootstrapped")
	}

	// Create collections and indexes
	db := RootDb()
	log.Printf("Creating collections and indexes...")
	db.C("config").Create(&mgo.CollectionInfo{})
	db.C("tenants").Create(&mgo.CollectionInfo{})
	db.C("tenants").EnsureIndex(mgo.Index{Key: []string{"code"}, Unique: true})
	db.C("users").Create(&mgo.CollectionInfo{})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"apikey"}, Unique: true})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"tenantid"}, Unique: true})

	// Create default tenant
	/*
		tenant := models.Tenant{
			Name: answers.Tenant.Name,
		}
		tenant.Init(c.NewId())
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
		user.Init(c.NewId())
		user.TenantId = tenant.Id

		err = db.C("users").Insert(user)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created root user '%s %s <%s>' with ID %s", user.FirstName, user.LastName, user.Email, user.Id.(bson.ObjectId).Hex())
	*/
	// Create common.entry
	apiInfo := M{
		"Version":     "1.0.0",
		"installDate": time.Now(),
	}
	err = db.C("config").Insert(apiInfo)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal("Configuration initialization completed successfully")

	return nil
}

func NewId() interface{} {
	return bson.NewObjectId()
}

func IdToString(id interface{}) string {
	oid, ok := id.(bson.ObjectId)

	if ok {
		return oid.Hex()
	}

	panic(fmt.Sprintf("Unknown ID format (%s)", reflect.TypeOf(id)))
}

func CreateDatabase(database string) error {
	return nil
}

func DeleteDatabase(database string) error {
	err := dbSession.DB(database).DropDatabase()

	return err
}
