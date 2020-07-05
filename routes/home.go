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
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

// HomepageHandler response for the home page.
func HomepageHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	courses := models.GetCourses()
	for i := range courses {
		courses[i].LoadPostsCount()
	}
	ctx.Data["Courses"] = courses
	ctx.Data["Home"] = 1
	ctx.HTML(200, "index")
}

// QnAHandler response for the Questions and Answers page.
func QnAHandler(ctx *macaron.Context, sess session.Store) {
	ctx.Data["Title"] = "Q&A"
	ctx.HTML(200, "qna")
}

// GuidelinesHandler response for the Guidelines page.
func GuidelinesHandler(ctx *macaron.Context, sess session.Store) {
	ctx.Data["Title"] = "Guidelines"
	ctx.HTML(200, "guidelines")
}
