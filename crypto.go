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
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
)

func GenerateApiKey(user User) string {
	// Create API key
	jsonHash, err := json.Marshal(&user)
	if err != nil {
		panic(err)
	}

	shaSum := sha1.Sum(jsonHash)
	return fmt.Sprintf("%x", shaSum)
}

func RandomSalt() []byte {
	// Generate a random salt
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		log.Panic(err)
	}

	return salt
}

func HashPasswordWithSalt(password string, salt []byte) string {
	// Prepend the salt with the password
	salted := append(salt, []byte(password)...)

	// Hash it up
	sha := sha256.Sum256(salted)

	// Store the salt for later
	hash := append(sha[:], salt...)

	// Base64 encode
	return base64.StdEncoding.EncodeToString(hash[:])
}

func HashPassword(password string) string {
	salt := RandomSalt()
	return HashPasswordWithSalt(password, salt)
}

func CheckPassword(hash string, password string) bool {
	// Decode base64 hash to [32]byte SHA256 sum
	b, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Panic(err)
	}

	// Compare
	return hash == HashPasswordWithSalt(password, b[32:])
}
