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
 */
package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type TimeStampFormat struct{}

func (c *TimeStampFormat) GetName() string {
	return "timestamp"
}

func (c *TimeStampFormat) Validate(att *CITypeAttribute, val *interface{}) error {
	if att.Type != c.GetName() {
		return errors.New(fmt.Sprintf("Attribute '%s' is not the correct type", att.Name))
	}

	// Timestamps to be stored as Int64 of milliseconds since 1970-01-01T00:00:00.000Z
	if _, ok := (*val).(float64); ok {
		return nil
	}

	// Parse string formats
	if str, ok := (*val).(string); ok {

		// Milliseconds since 1970
		i64, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			(*val) = i64
			return nil
		}

		// Convert to a time.Time first
		layouts := []string{
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.RFC1123,
			time.RFC1123Z,
			time.RFC3339,
		}

		for _, layout := range layouts {
			t, err := time.Parse(layout, str)
			if err == nil {
				// Convert to milliseconds since 1970
				(*val) = (int64(1000) * t.Unix()) + int64(1000000*t.Nanosecond())
				return nil
			}
		}
	}

	return errors.New(fmt.Sprintf("Value '%v' for attribute '%s' is not a valid timestamp", *val, att.Name))
}
