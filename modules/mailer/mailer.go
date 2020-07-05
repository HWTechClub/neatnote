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
package mailer

import (
	"git.sr.ht/~humaid/neatnote/modules/settings"
	"net/smtp"
	"strings"
)

// EmailCode emails the provided code to the provided email address.
func EmailCode(to string, code string) (err error) {
	from := settings.Config.EmailAddress
	message := "From: <" + from + ">\n" +
		"To: <" + to + ">\n" +
		"Subject: " + settings.Config.SiteName + " login code\n\n" +
		"Hello!\nYour login code is " + code + "\n\n" +
		"Ignore this message if you have not requested a login.\n\n\n" +
		"- " + settings.Config.SiteName + "\nThis message is sent from an unmonitored inbox."

	err = smtp.SendMail(settings.Config.EmailSMTPServer,
		smtp.PlainAuth("", from, settings.Config.EmailPassword, strings.Split(settings.Config.EmailSMTPServer, ":")[0]),
		from, []string{to}, []byte(message))

	return err
}
