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
	"git.sr.ht/~humaid/neatnote/modules/common"
)

// Post represents a post in one of the Courses by one of the Users.
// It keeps track of the poster, course, comments (and comment count), whether
// it is an anonymous post, and so on.
type Post struct {
	PostID        int64     `xorm:"pk autoincr"`
	CourseCode    string    `xorm:"text notnull"`
	PosterID      string    `xorm:"notnull"`
	Poster        *User     `xorm:"-" json:"-"`
	Locked        bool      `xorm:"notnull"` // Whether the comments are locked.
	Comments      []Comment `xorm:"-" json:"-"`
	CommentsCount int64     `xorm:"-" json:"-"`
	Title         string    `xorm:"text notnull"`
	Text          string    `xorm:"text notnull"`
	CreatedUnix   int64     `xorm:"created"`
	UpdatedUnix   int64     `xorm:"updated"`
	Anonymous     bool      `xorm:"notnull"`
	AnonName      string    `xorm:"text null"`
	Iota          int64
}

// LoadComments loads the comments of the post into a non-mapped field.
func (p *Post) LoadComments() (err error) {
	return engine.Where("post_id = ?", p.PostID).Find(&p.Comments)
}

// GetPost gets a Post based on the ID.
// It will return the pointer to the Post, and whether there was an error.
func GetPost(id string) (*Post, error) {
	p := new(Post)
	has, err := engine.ID(id).Get(p)
	if err != nil {
		return p, err
	} else if !has {
		return p, errors.New("Post does not exist")
	}
	p.Poster, _ = GetUser(p.PosterID)
	return p, nil
}

// UpdatePost updates a post in the database.
func UpdatePost(p *Post) (err error) {
	_, err = engine.Id(p.PostID).Update(p)
	return
}

// AddPost adds a new Post to the database.
func AddPost(p *Post) (err error) {
	_, err = engine.Insert(p)
	return err
}

// DeletePost deletes a post from the database, including comments.
func DeletePost(id string) (err error) {
	sess := engine.NewSession()

	_, err = sess.Id(id).Delete(&Post{})
	if err != nil {
		return err
	}
	_, err = sess.Where("post_id = ?", id).Delete(new(Comment))
	if err != nil {
		return err
	}

	return sess.Commit()
}

// GetAllUserPosts returns all of the post created by a specific user.
func GetAllUserPosts(user string) (p []Post, err error) {
	err = engine.Where("poster_id = ?", user).Find(&p)
	return
}

// GetUserPostCount returns number of posts by a specific user.
func GetUserPostCount(user string) (i int64, err error) {
	return engine.Where("poster_id = ?", user).Count(new(Post))
}

func UnvotePost(user string, post int64) (err error) {
	sess := engine.NewSession()
	if err = sess.Begin(); err != nil {
		return err
	}
	u := new(User)
	var has bool
	if has, err = sess.ID(user).Get(u); err != nil {
		return err
	} else if !has {
		return errors.New("User does not exist.")
	}

	if !common.ContainsInt64(u.Upvoted, post) {
		return errors.New("Cannot unvote a post which is not upvoted.")
	}

	p := new(Post)
	has, err = sess.ID(post).Get(p)
	if has, err = sess.ID(user).Get(u); err != nil {
		return err
	} else if !has {
		return errors.New("Post does not exist.")
	}

	_, err = sess.Id(u.Username).Cols("upvoted").Update(&User{
		Username: u.Username,
		Upvoted:  common.RemoveInt64(u.Upvoted, post),
	})
	if err != nil {
		return err
	}

	_, err = sess.ID(post).Cols("iota").Update(&Post{
		PostID: post,
		Iota:   p.Iota - 1,
	})
	if err != nil {
		return err
	}

	return sess.Commit()
}

func UpvotePost(user string, post int64) (err error) {
	sess := engine.NewSession()
	if err = sess.Begin(); err != nil {
		return err
	}
	u := new(User)
	var has bool
	if has, err = sess.ID(user).Get(u); err != nil {
		return err
	} else if !has {
		return errors.New("User does not exist.")
	}

	if common.ContainsInt64(u.Upvoted, post) {
		return errors.New("Already upvoted")
	}

	p := new(Post)
	has, err = sess.ID(post).Get(p)
	if has, err = sess.ID(user).Get(u); err != nil {
		return err
	} else if !has {
		return errors.New("Post does not exist.")
	}

	_, err = sess.Id(u.Username).Cols("upvoted").Update(&User{
		Username: u.Username,
		Upvoted:  append(u.Upvoted, post),
	})
	if err != nil {
		return err
	}

	_, err = sess.ID(post).Cols("iota").Update(&Post{
		PostID: post,
		Iota:   p.Iota + 1,
	})
	if err != nil {
		return err
	}

	return sess.Commit()
}
