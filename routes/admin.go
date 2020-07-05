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
	"git.sr.ht/~humaid/neatnote/models"
	"git.sr.ht/~humaid/neatnote/modules/settings"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/hako/durafmt"
	macaron "gopkg.in/macaron.v1"
	"runtime"
	"time"
)

// AdminHandler response for the admin dashboard.
func AdminHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	ctx.Data["Title"] = "Admin Dashboard"
	ctx.Data["Goroutines"] = runtime.NumGoroutine()
	ctx.Data["Goversion"] = runtime.Version()
	ctx.Data["Uptime"] = durafmt.Parse(time.Now().Sub(settings.StartTime)).String()
	ctx.Data["Users"] = models.GetUsers()
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "admin/index")
}

// AdminViewUserHandler response for viewing a user's information.
func AdminViewUserHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	u, err := models.GetUser(ctx.Params("user"))
	if err != nil {
		f.Error("User not found.")
		ctx.Redirect("/a")
		return
	}
	ctx.Data["VUser"] = u
	ctx.HTML(200, "admin/view-user")

}

// PostAdminHandler post response for admin dashboard.
func PostAdminHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	ctx.Data["csrf_token"] = x.GetToken()
	switch ctx.Query("action") {
	case "suspend_prompt":
		u, err := models.GetUser(ctx.Query("username"))
		if err != nil {
			panic(err)
		}
		if u.IsAdmin {
			f.Error("Cannot suspend an admin.")
			break
		}
		ctx.Data["BadUser"] = u
		ctx.HTML(200, "admin/suspend")
		return
	case "suspend":
		u, err := models.GetUser(ctx.Query("username"))
		if err != nil {
			panic(err)
		}
		if u.IsAdmin {
			f.Error("Cannot suspend an admin.")
			break
		}
		if len(ctx.QueryTrim("reason")) < 3 {
			f.Error("You have to specify a reason.")
			ctx.Data["BadUser"] = u
			ctx.HTML(200, "admin/suspend")
			return
		}
		u.Suspended = true
		u.SuspendReason = ctx.QueryTrim("reason")
		models.UpdateUserCols(u, "suspended", "suspend_reason")
	case "unsuspend":
		u, err := models.GetUser(ctx.Query("username"))
		if err != nil {
			panic(err)
		}
		if !u.Suspended {
			f.Error("User already unsuspended!")
		}
		u.Suspended = false
		models.UpdateUserCols(u, "suspended")
	default:
		f.Error("Unknown action.")
	}
	ctx.Redirect("/a")
}

// AdminAddCourseHandler response for adding a new course.
func AdminAddCourseHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.Data["Title"] = "Add course"
	ctx.HTML(200, "admin/add-course")
}

// AdminPostAddCourseHandler post response for adding a new course.
func AdminPostAddCourseHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	courseCode := ctx.QueryTrim("coursecode")
	courseName := ctx.QueryTrim("coursename")

	// Check if course exists already
	if len(courseCode) < 1 || len(courseName) < 1 {
		f.Error("You must specify course code and name!")
		ctx.Redirect("/a/addcourse")
		return
	} else if _, err1 := models.GetCourse(courseCode); err1 == nil {
		f.Error("Course already exists!")
		ctx.Redirect("/a/addcourse")
		return
	}

	models.AddCourse(&models.Course{
		Code:    courseCode,
		Name:    courseName,
		Visible: true,
		Locked:  false,
	})

	ctx.Redirect("/")
}
