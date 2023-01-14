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

package main

//defines misc functions
import (
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func create_bind(user string, channel string) {
	if find_bind(user) != "" {
		return
	}
	bind := binding{User: user, Channel: channel}
	bindings = append(bindings, bind)
}

func find_bind(user string) string {
	for x := 0; x < len(bindings); x++ {
		if bindings[x].User == user {
			return bindings[x].Channel
		}
	}
	return ""
}

func modify_bind(user string, channel string) bool {
	for x := 0; x < len(bindings); x++ {
		if bindings[x].User == user {
			bindings[x].Channel = channel
			return true
		}
	}
	return false
}
