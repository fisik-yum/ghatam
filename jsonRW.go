/*
   kavaca- a discord bot that acts as a mail forwarder
   Copyright (C) 2021  fisik_yum

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// read (and write) configuration files
package main

import (
	"encoding/json"
	"log"
	"os"
)

func read_config() { // main config file for end user
	f, err := os.ReadFile("config.json")
	check(err)
	var userData owner
	err = json.Unmarshal([]byte(f), &userData)
	check(err)
	token = userData.Token
	ownerID = userData.ID
	defaultChannel = userData.DefaultChannel
	if token == "" || ownerID == "" {
		panic("Required config info missing")
	}
}

func save_bindings() { //add sanitize
	fd, err := json.Marshal(bindings)
	check(err)
	err = os.WriteFile("bindings.json", fd, 0777)
	check(err)
}

func load_bindings() {
	bindings = nil
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic:", err)
		}
	}()
	f, err := os.ReadFile("bindings.json")
	check(err)
	err = json.Unmarshal([]byte(f), &bindings)
	check(err)
}
