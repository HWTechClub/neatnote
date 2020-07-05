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
	"git.sr.ht/~humaid/neatnote/modules/mailer"
	"git.sr.ht/~humaid/neatnote/modules/settings"
	"github.com/badoux/checkmail"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
	"math/rand"
	"strings"
)

// LogoutHandler response for logging out.
func LogoutHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	sess.Set("auth", LoggedOut)
	//sess.Flush()
	ctx.Redirect("/")
}

// LoginHandler response for logging in page.
func LoginHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	if sess.Get("auth") == Verification {
		ctx.Redirect("/verify")
		return
	} else if sess.Get("auth") == LoggedIn {
		f.Info("You are already logged in!")
		ctx.Redirect("/")
		return
	}
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.Data["Title"] = "Login"
	ctx.HTML(200, "login")
}

func randIntRange(min, max int) int {
	return rand.Intn(max-min) + min
}

// PostLoginHandler post response for login page.
func PostLoginHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	if sess.Get("auth") == Verification {
		f.Warning("You need to verify before you continue.")
		ctx.Redirect("/verify")
		return
	} else if sess.Get("auth") == LoggedIn {
		f.Info("You are already logged in!")
		ctx.Redirect("/")
		return
	}
	// Generate code
	code := fmt.Sprint(randIntRange(100000, 999999))
	to := fmt.Sprintf("%s%s", strings.ToLower(ctx.QueryTrim("email")), settings.Config.UniEmailDomain)
	err := checkmail.ValidateFormat(to)
	if err != nil {
		f.Error("You provided an invalid email.")
		ctx.Redirect("/login")
		return
	}

	if !settings.Config.DevMode {
		go mailer.EmailCode(to, code)
	}
	sess.Set("auth", Verification)
	sess.Set("code", code)
	sess.Set("user", to)
	sess.Set("attempts", 0)
	ctx.Redirect("/verify")
}

// VerifyHandler response for verification page.
func VerifyHandler(ctx *macaron.Context, x csrf.CSRF, sess session.Store, f *session.Flash) {
	if sess.Get("auth") == LoggedOut {
		ctx.Redirect("/login")
		return
	} else if sess.Get("auth") == LoggedIn || sess.Get("auth") != Verification {
		f.Info("You are already logged in!")
		ctx.Redirect("/")
		return
	}
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.Data["email"] = sess.Get("user")
	ctx.Data["Title"] = "Verification"
	ctx.HTML(200, "verify_login")
}

// CancelHandler post response for canceling verification.
func CancelHandler(ctx *macaron.Context, sess session.Store) {
	if sess.Get("auth") != Verification {
		ctx.Redirect("/login")
		return
	}

	sess.Set("auth", LoggedOut)
	ctx.Redirect("/login")
}

// PostVerifyHandler post reponse for verification.
func PostVerifyHandler(ctx *macaron.Context, sess session.Store, f *session.Flash) {
	if sess.Get("auth") == LoggedOut {
		ctx.Redirect("/login")
		return
	} else if sess.Get("auth") == LoggedIn {
		ctx.Redirect("/")
		return
	}
	sess.Set("attempts", sess.Get("attempts").(int)+1)
	if sess.Get("attempts").(int) > 3 {
		f.Error("You reached the maximum number of attempts. Please try again later.")
		sess.Set("auth", LoggedOut)
		ctx.Redirect("/")
		return
	}
	if ctx.QueryTrim("code") != sess.Get("code") && !settings.Config.DevMode {
		f.Error("The code you entered is invalid, make sure you use the latest code sent to you.")
		ctx.Redirect("/verify")
		return
	}
	if !models.HasUser(sess.Get("user").(string)) {
		models.AddUser(&models.User{
			Username: sess.Get("user").(string),
			IsAdmin:  false,
		})
	} else {
		u, err := models.GetUser(sess.Get("user").(string))
		if err != nil {
			panic(err)
		}
		if u.Suspended {
			sess.Set("auth", LoggedOut)
			ctx.Data["SuspendReason"] = u.SuspendReason
			ctx.HTML(200, "suspended")
			return
		}
	}

	sess.Set("auth", LoggedIn)
	ctx.Redirect("/")
}
