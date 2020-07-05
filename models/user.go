// Neat Note. A notes sharing platform for university students.
// Copyright (C) 2020 Humaid AlQassimi
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package models

import (
	"errors"
)

// User represents a website user.
// It keeps track of the iota, settings (such as badges), and whether they
// have administrative privileges.
type User struct {
	Username      string `xorm:"pk"`
	FullName      string `xorm:"text null"`
	Badge         string `xorm:"text null"`
	IsAdmin       bool   `xorm:"bool"`
	Iota          int64
	CreatedUnix   int64   `xorm:"created"`
	Upvoted       []int64 // Post IDs which the user upvoted.
	Suspended     bool    `xorm:"notnull"`
	SuspendReason string  `xorm:"text null"`
}

// GetUser gets a user based on their username.
func GetUser(user string) (*User, error) {
	u := new(User)
	has, err := engine.ID(user).Get(u)
	if err != nil {
		return u, err
	} else if !has {
		return u, errors.New("User does not exist")
	}
	u.Iota, _ = engine.Where("poster_id = ?", u.Username).SumInt(new(Post), "iota")
	return u, nil
}

// GetUsers returns a list of all users in the database.
func GetUsers() (users []User) {
	engine.Find(&users)
	return users
}

// AddUser adds a new User to the database.
func AddUser(u *User) (err error) {
	_, err = engine.Insert(u)
	return err
}

// HasUser returns whether a user exists in the database.
func HasUser(user string) (has bool) {
	has, _ = engine.Get(&User{Username: user})
	return has
}

// UpdateUser updates a user in the database.
func UpdateUser(u *User) (err error) {
	_, err = engine.Id(u.Username).Update(u)
	return
}

// UpdatUserCols updates a user in the database including the specified
// columns, even if the fields are empty.
func UpdateUserCols(u *User, cols ...string) error {
	_, err := engine.Id(u.Username).Cols(cols...).Update(u)
	return err
}
