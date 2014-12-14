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
package services

import (
        "log"
        
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
)

type Database struct {
    *mgo.Session
}

// Wire the service
func DatabaseService(session *mgo.Session) martini.Handler {
    if session == nil {
        var err error
        session, err = mgo.Dial("mongodb://localhost")
        if err != nil {
                log.Panic(err)
        }
    }
    
    return func(c martini.Context) {
	newSession := session.Clone()
	db := Database{newSession}
        c.Map(&db)
        defer newSession.Close()
        c.Next()
    }
}
