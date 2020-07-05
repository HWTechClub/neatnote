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
package routes

import (
	"fmt"
	"git.sr.ht/~humaid/neatnote/models"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

// CourseExists is a per-route middleware which checks if the course exists
// in the database, otherwise display an error.
func CourseExists(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	c, err := models.GetCourse(ctx.Params("course"))
	if err != nil {
		f.Error("Welp! The course doesn't exist.")
		ctx.Redirect("/")
		return
	}
	ctx.Data["Course"] = c
}

// PostExists is a per-route middleware which checks if the course and post
// exists in the database, otherwise display an error.
func PostExists(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	c, err := models.GetCourse(ctx.Params("course"))
	if err != nil {
		f.Error("Welp! The course doesn't exist.")
		ctx.Redirect("/")
		return
	}
	ctx.Data["Course"] = c

	p, err := models.GetPost(ctx.Params("post"))
	if err != nil {
		f.Error("Post does not exist.")
		ctx.Redirect(fmt.Sprintf("/c/%s", ctx.Params("course")))
		return
	}
	if p.CourseCode != c.Code {
		ctx.Redirect(fmt.Sprintf("/c/%s/%d", p.CourseCode, p.PostID))
		return
	}
	ctx.Data["Post"] = p

}

// CourseUnlocked is a per-route middleware which checks if the course is
// unlocked, otherwise display an error.
func CourseUnlocked(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	if c, err := models.GetCourse(ctx.Params("course")); err != nil {
		panic(err)
	} else if c.Locked {
		f.Error("This course is locked.")
		ctx.Redirect(fmt.Sprintf("/c/%s", ctx.Params("course")))
		return
	}
}

// PostUnlocked is a per-route middleware which checks if the course and post
// is unlocked, otherwise display an error.
func PostUnlocked(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	if c, err := models.GetCourse(ctx.Params("course")); err != nil {
		panic(err)
	} else if c.Locked {
		f.Error("You cannot do that as the course is locked.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}
	if post, err := models.GetPost(ctx.Params("post")); err != nil {
		panic(err)
	} else if post.Locked {
		f.Error("You cannot do that as the post is locked.")
		ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
			ctx.Params("post")))
		return
	}
}

// RequireLogin is a per-route middleware which checks if the user is logged
// in, othwerise redirect to the login page and display an error.
func RequireLogin(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	if sess.Get("auth") != LoggedIn {
		f.Error("Please login first!")
		ctx.Redirect("/login")
		return
	}
}

// RequireAdmin is a per-route middleware which checks if the user is an admin,
// otherwise display an error.
func RequireAdmin(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	u, err := models.GetUser(sess.Get("user").(string))
	if err != nil {
		panic(err)
	}
	if !u.IsAdmin {
		f.Error("You may not do that.")
		if len(ctx.Params("post")) > 0 {
			ctx.Redirect(fmt.Sprintf("/c/%s/%s", ctx.Params("course"),
				ctx.Params("post")))
		} else {
			ctx.Redirect("/")
		}
		return
	}
}
