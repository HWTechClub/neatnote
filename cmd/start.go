package cmd

import (
	"fmt"
	"git.sr.ht/~humaid/notes-overflow/models"
	"git.sr.ht/~humaid/notes-overflow/modules/settings"
	"git.sr.ht/~humaid/notes-overflow/routes"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	_ "github.com/go-macaron/session/postgres"
	"github.com/urfave/cli"
	macaron "gopkg.in/macaron.v1"
	"log"
	"net/http"
)

// CmdStart represents a command-line command
// which starts the bot.
var CmdStart = cli.Command{
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

	m.Use(macaron.Renderer())
	m.Use(cache.Cacher())
	psqlConfig := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s sslmode=disable",
		settings.DBConfig.User, settings.DBConfig.Password, settings.DBConfig.Host, "session")
	fmt.Println(psqlConfig)
	m.Use(session.Sessioner(session.Options{
		Provider:       "postgres",
		ProviderConfig: psqlConfig,
	}))
	m.Use(csrf.Csrfer())
	m.Use(captcha.Captchaer())

	// Web routes
	m.Get("/", routes.HomepageHandler)

	// Login and verification
	m.Get("/login", routes.LoginHandler)
	m.Post("/login", routes.PostLoginHandler)
	m.Get("/logout", routes.LogoutHandler)
	m.Get("/verify", routes.VerifyHandler)
	m.Post("/verify", routes.PostVerifyHandler)

	log.Printf("Starting web server on port %s\n", settings.SitePort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", settings.SitePort), m))
	return nil
}
