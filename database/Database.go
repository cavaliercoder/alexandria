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
package database

import (
	"github.com/cavaliercoder/alexandria/common"
	"github.com/cavaliercoder/alexandria/models"
	"fmt"
)

type Driver interface {
	Connect(*common.DatabaseConfig) error              // Connect to the database
	Clone() (Driver, error)                            // Clone a database connection
	Close() error                                      // Disconnect from the database
	Use(string) error                                  // Select a database for use
	IsBootStrapped() (bool, error)                     // Return true is datasbe schema is intialized
	BootStrap(*common.Answers) error                   // Initialize database schema
	CreateDatabase(string) error                       // Create a new CMDB
	DeleteDatabase(string) error                       // Delete a CMDB
	NewId() interface{}				   // Create a new unique record ID
	IdToString(interface{}) string                     // Convert a database ID record to string
	GetAll(string, M, interface{}) error               // Get multiple entities from the database
	GetOne(string, M, interface{}) error               // Get a single entity from the database
	GetOneById(string, interface{}, interface{}) error // Get a single entity from the database by ID
	Insert(string, models.Model) error                 // Add an entity to the database
	Delete(string, M) error                            // Delete a resource from the CMDB
}

// M is a convenience shortcut for `map[string]interface{}`
type M map[string]interface{}

/*
 * Connect loads application common.to select a database driver, loads
 * the driver and connects to the database using the specified connection
 * common.
 */
func Connect() (Driver, error) {
	var driver Driver
	config, err := common.GetConfig()
	if err != nil {
		return nil, err
	}

	switch config.Database.Driver {
	case "mongodb":
		// Connect to database
		driver = new(MongoDriver)
		err = driver.Connect(&config.Database)
		if err != nil {
			return nil, err
		}

	default:
		panic(fmt.Sprintf("Unsupported database driver: %s", config.Database.Driver))
	}

	return driver, nil
}
