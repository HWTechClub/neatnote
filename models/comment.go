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
	"html/template"
)

// Comment represents a comment on a Post. It keeps track of the poster and
// which post it is posted to.
type Comment struct {
	CommentID     int64         `xorm:"pk autoincr"`
	PostID        int64         `xorm:"notnull"`
	PosterID      string        `xorm:"notnull"`
	Poster        *User         `xorm:"-" json:"-"`
	Text          string        `xorm:"notnull"`
	FormattedText template.HTML `xorm:"-" json:"-"`
	CreatedUnix   int64         `xorm:"created"`
	UpdatedUnix   int64         `xorm:"updated"`
}

// LoadPoster loads the poster of a comment in the non-mapped field of the
// Comment struct.
func (c *Comment) LoadPoster() (err error) {
	if c == nil {
		return nil
	} else if c.Poster != nil {
		return nil
	}

	c.Poster, err = GetUser(c.PosterID)
	return
}

// AddComment adds a new Comment to the database.
func AddComment(c *Comment) (err error) {
	_, err = engine.Insert(c)
	return err
}

// UpdateComment updates a comment in the database.
func UpdateComment(c *Comment) (err error) {
	_, err = engine.Id(c.CommentID).Update(c)
	return
}

// GetComment gets a comment based on the ID.
// It will return the pointer to the Comment, and whether there was an error.
func GetComment(id string) (*Comment, error) {
	c := new(Comment)
	has, err := engine.ID(id).Get(c)
	if err != nil {
		return c, err
	} else if !has {
		return c, errors.New("Comment does not exist")
	}
	return c, nil
}

// DeleteComment deletes a comment from the database.
func DeleteComment(id string) (err error) {
	_, err = engine.Id(id).Delete(&Comment{})
	return
}

// GetAllUserComments returns all of the comments created by a specific user.
func GetAllUserComments(user string) (c []Comment, err error) {
	err = engine.Where("poster_id = ?", user).Find(&c)
	return
}
