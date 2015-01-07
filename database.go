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
 */
package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
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
	count, err := RootDb().C("apiInfo").Find(nil).Count()
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

	config, err := GetConfig()
	if err != nil {
		return err
	}

	// Create collections and indexes
	db := RootDb()
	log.Printf("Creating collections and indexes...")
	db.C("apiInfo").Create(&mgo.CollectionInfo{})

	db.C("tenants").Create(&mgo.CollectionInfo{})
	db.C("tenants").EnsureIndex(mgo.Index{Key: []string{"code"}, Unique: true})

	db.C("users").Create(&mgo.CollectionInfo{})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"email"}, Unique: true})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"apikey"}, Unique: true})
	db.C("users").EnsureIndex(mgo.Index{Key: []string{"tenantid"}, Unique: false})

	// Create default tenant
	tenant := Tenant{
		Name: answers.Tenant.Name,
	}
	tenant.InitModel()
	err = db.C("tenants").Insert(tenant)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created detault tenant '%s' with code %s", tenant.Name, tenant.Code)

	// Create root user
	user := User{
		FirstName:    answers.User.FirstName,
		LastName:     answers.User.LastName,
		Email:        answers.User.Email,
		PasswordHash: HashPassword(answers.User.Password),
	}
	user.InitModel()
	user.TenantId = tenant.Id

	// Preset ApiKey for dev
	if !config.Server.Production {
		user.ApiKey = "D8fzx4cpX0SrPm6cEb6HwLf6IvCb0MvA"
	}

	err = db.C("users").Insert(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created root user '%s %s <%s>'", user.FirstName, user.LastName, user.Email)

	// Create config entry
	apiInfo := ApiInfo{
		Version:     "1.0.0",
		InstallDate: time.Now(),
		RootUserId:  user.Id,
	}
	err = db.C("apiInfo").Insert(apiInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Write config to rc file
	rcfile := ExpandPath("~/.alexrc")
	file, err := os.Create(rcfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("ALEX_API_URL=\"http://localhost:%d%s\"\nALEX_API_KEY=\"%s\"\nALEX_API_DB=\"%s\"\n", config.Server.ListenPort, ApiV1Prefix, user.ApiKey, config.Database.Database))
	file.Sync()
	log.Printf("Saved Alexandria CMDB configuration to %s", rcfile)

	log.Print("Configuration initialization completed successfully")
	os.Exit(0)

	return nil
}

func CreateCmdb(name string) error {
	session := DbConnect()
	db := session.DB(name)

	// Create CI Types collection
	err := db.C("citypes").Create(&mgo.CollectionInfo{})
	if err != nil {
		return err
	}

	err = db.C("citypes").EnsureIndex(mgo.Index{Key: []string{"shortname"}, Unique: true})
	if err != nil {
		return err
	}

	return err
}

func DropCmdb(name string) error {
	session := DbConnect()
	db := session.DB(name)
	err := db.DropDatabase()

	return err
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

func IdFromString(id string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(id) {
		return "", errors.New(fmt.Sprintf("Invalid ID: %s", id))
	}

	return bson.ObjectIdHex(id), nil
}

func CreateDatabase(database string) error {
	return nil
}

func DeleteDatabase(database string) error {
	err := dbSession.DB(database).DropDatabase()

	return err
}
