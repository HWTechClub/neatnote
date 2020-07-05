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
package cmd

import (
	"fmt"
	"git.sr.ht/~humaid/neatnote/models"
	"git.sr.ht/~humaid/neatnote/modules/settings"
	"git.sr.ht/~humaid/neatnote/routes"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	_ "github.com/go-macaron/session/mysql" // MySQL driver for persistent sessions
	"github.com/hako/durafmt"
	"github.com/urfave/cli/v2"
	macaron "gopkg.in/macaron.v1"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// CmdStart represents a command-line command
// which starts the bot.
var CmdStart = &cli.Command{
	Name:    "run",
	Aliases: []string{"start", "web"},
	Usage:   "Start the web server",
	Action:  start,
}

func start(clx *cli.Context) (err error) {
	settings.LoadConfig()
	engine := models.SetupEngine()
	defer engine.Close()

	// Run macaron
	m := macaron.Classic()
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{map[string]interface{}{
			"CalcTime": func(sTime time.Time) string {
				return fmt.Sprint(time.Since(sTime).Nanoseconds() / int64(time.Millisecond))
			},
			"EmailToUser": func(s string) string {
				if strings.Contains(s, "@") {
					return strings.Split(s, "@")[0]
				} else {
					return s
				}
			},
			"CalcDurationShort": func(unix int64) string {
				return durafmt.Parse(time.Now().Sub(time.Unix(unix, 0))).LimitFirstN(1).String()
			},
		}},
		IndentJSON: true,
	}))

	if settings.Config.DevMode {
		fmt.Println("In development mode.")
		macaron.Env = macaron.DEV
	} else {
		fmt.Println("In production mode.")
		macaron.Env = macaron.PROD
	}

	m.Use(cache.Cacher())
	sessOpt := session.Options{
		CookieLifeTime: 15778800, // 6 months
		Gclifetime:     15778800,
		CookieName:     "hithereimacookie",
	}
	if settings.Config.DBConfig.Type == settings.MySQL {
		sqlConfig := fmt.Sprintf("%s:%s@tcp(%s)/%s",
			settings.Config.DBConfig.User, settings.Config.DBConfig.Password,
			settings.Config.DBConfig.Host, settings.Config.DBConfig.Name)
		sessOpt.Provider = "mysql"
		sessOpt.ProviderConfig = sqlConfig
		sessOpt.CookieLifeTime = 0
	}
	m.Use(session.Sessioner(sessOpt))
	m.Use(csrf.Csrfer())
	m.Use(captcha.Captchaer())
	m.Use(routes.ContextInit())

	// Web routes
	m.Get("/", routes.HomepageHandler)
	m.Group("/profile", func() {
		m.Get("/", routes.ProfileHandler)
		m.Post("/", csrf.Validate, routes.PostProfileHandler)
		m.Post("/data.json", csrf.Validate, routes.PostDataHandler)
	}, routes.RequireLogin)
	m.Get("/qna", routes.QnAHandler)
	m.Get("/guidelines", routes.GuidelinesHandler)

	// Login and verification
	m.Get("/login", routes.LoginHandler)
	m.Post("/login", csrf.Validate, routes.PostLoginHandler)
	m.Get("/verify", routes.VerifyHandler)
	m.Post("/verify", csrf.Validate, routes.PostVerifyHandler)
	m.Post("/cancel", csrf.Validate, routes.CancelHandler)
	m.Get("/logout", routes.RequireLogin, routes.LogoutHandler)

	m.Group("/a", func() {
		m.Get("/", routes.AdminHandler)
		m.Post("/", csrf.Validate, routes.PostAdminHandler)
		m.Get("/view/:user", routes.AdminViewUserHandler)
		m.Get("/addcourse", routes.AdminAddCourseHandler)
		m.Post("/addcourse", csrf.Validate, routes.AdminPostAddCourseHandler)
	}, routes.RequireLogin, routes.RequireAdmin)

	m.Group("/c/:course", func() {
		m.Get("/", routes.CourseHandler)
		m.Get("/:sort(new|top)", routes.CourseHandler)
		m.Get("/post", routes.RequireLogin, routes.CourseUnlocked,
			routes.CreatePostHandler)
		m.Post("/post", routes.RequireLogin, routes.CourseUnlocked,
			csrf.Validate, routes.PostCreatePostHandler)
		m.Group("/:post", func() {
			m.Get("/", routes.PostPageHandler)
			m.Get("/lite", routes.LitePostHandler)
			m.Post("/", csrf.Validate, routes.PostCommentPostHandler)
			m.Get("/upvote", routes.RequireLogin, routes.PostUnlocked,
				routes.UpvotePostHandler)
			m.Get("/edit", routes.RequireLogin, routes.PostUnlocked,
				routes.EditPostHandler)
			m.Post("/edit", routes.RequireLogin, routes.PostUnlocked,
				routes.PostEditPostHandler)
			m.Get("/edit/:id", routes.RequireLogin, routes.PostUnlocked,
				routes.EditCommentHandler)
			m.Post("/edit/:id", routes.RequireLogin, routes.PostUnlocked,
				routes.PostEditCommentHandler)
			m.Group("/", func() {
				m.Get("/del/:id", routes.PostUnlocked,
					routes.DeleteCommentHandler)
				m.Get("/del", routes.PostUnlocked, routes.DeletePostHandler)
				m.Get("/reveal", routes.RevealPosterHandler)
			}, routes.RequireLogin, routes.RequireAdmin)
		}, routes.PostExists)
	}, routes.CourseExists)

	log.Printf("Starting web server on port %s\n", settings.Config.SitePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", settings.Config.SitePort), m))
	return nil
}
