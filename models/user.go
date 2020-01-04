package models

import (
	"errors"
)

// User represents a website user.
// It keeps track of the iota, settings (such as badges), and whether they
// have administrative privileges.
type User struct {
	Username    string `xorm:"pk"`
	FullName    string `xorm:"text null"`
	Badge       string `xorm:"text null"`
	IsAdmin     bool   `xorm:"bool"`
	Iota        int64
	Created     string  `xorm:"-"`
	CreatedUnix int64   `xorm:"created"`
	Upvoted     []int64 // Post IDs which the user upvoted.
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
	u.Created = calcDuration(u.CreatedUnix)
	return u, nil
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

// UpdateUserBadge updates a user in the database including the Badge field,
// even if the field is empty.
func UpdateUserBadge(u *User) (err error) {
	_, err = engine.Id(u.Username).Cols("badge").Update(u)
	return
}