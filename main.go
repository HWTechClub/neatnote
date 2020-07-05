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
package main

import (
	"git.sr.ht/~humaid/neatnote/cmd"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// VERSION specifies the version of neatnote
var VERSION = "0.3.2"

func main() {
	app := &cli.App{
		Name:    "neatnote",
		Usage:   "a web app to allow University students to post notes in a civil manner.",
		Version: VERSION,
		Commands: []*cli.Command{
			cmd.CmdStart,
			cmd.CmdSetAdmin,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
