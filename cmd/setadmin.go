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
	"errors"
	"fmt"
	"git.sr.ht/~humaid/neatnote/models"
	"git.sr.ht/~humaid/neatnote/modules/settings"
	"github.com/urfave/cli/v2"
)

// CmdSetAdmin represents a command-line command
// which sets a user to become an admin.
var CmdSetAdmin = &cli.Command{
	Name:  "set-admin",
	Usage: "Gives a user admin status",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "status",
			Value: "true",
			Usage: "admin status",
		},
	},
	Action: setAdmin,
}

func setAdmin(c *cli.Context) error {
	if c.NArg() < 1 {
		return errors.New("Need username to set")
	}
	settings.LoadConfig()
	engine := models.SetupEngine()
	defer engine.Close()
	u, err := models.GetUser(c.Args().Get(0))
	if err != nil {
		return err
	}
	fmt.Printf("Current admin status: %t\n", u.IsAdmin)

	if c.String("status") == "true" {
		u.IsAdmin = true
		fmt.Println("Setting the user to admin")
	} else if c.String("status") == "false" {
		u.IsAdmin = false
		fmt.Println("Setting the user to non-admin")
	} else {
		return errors.New("Unknown status. Must be either true or false.")
	}

	return models.UpdateUserCols(u, "is_admin")
}
