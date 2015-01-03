/*
 * Alexandria CMDB - Open source config management database
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
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

type ServerConfig struct {
	Production bool   `json:"production"`
	ListenOn   string `json:"listenOn"`
	ListenPort int    `json:"listenPort"`
}

type DatabaseConfig struct {
	Driver   string   `json:"driver"`
	Servers  []string `json:"servers"`
	Timeout  int      `json:"timeout"`
	Database string   `json:"database"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

// default config file path
var confFilePath string = ""

// global, singleton config struct
var config *Config

func GetConfigFromFile(path string) (*Config, error) {
	if config != nil {
		return nil, errors.New("a configuration file was specified but configuration is already loaded")
	}

	confFilePath = path
	return GetConfig()
}

// GetConfig returns a pointer to a singleton config structure.
func GetConfig() (*Config, error) {
	if config == nil {
		// Select a config file
		if confFilePath == "" {
			if _, err := os.Stat("./config.json"); err == nil {
				confFilePath = "./config.json"
			} else if _, err := os.Stat("/etc/alexandria/config.json"); err == nil {
				confFilePath = "/etc/alexandria/config.json"
			} else {
				return nil, errors.New("no configuration file was found")
			}
		}

		// Open config file
		confFile, err := os.Open(confFilePath)
		if err != nil {
			return nil, err
		}

		defer confFile.Close()

		// Configuration defaults
		config = &Config{
			Database: DatabaseConfig{
				Driver:   "mongodb",
				Database: "alexandria",
			},
		}

		// Apply JSON config file
		parser := json.NewDecoder(confFile)
		if err = parser.Decode(config); err != nil {
			config = nil
			return nil, err
		}

		log.Printf("Loaded config from %s", confFilePath)
	}

	return config, nil
}
