package mailer

import (
	"git.sr.ht/~humaid/notes-overflow/modules/settings"
	"net/smtp"
)

func EmailCode(to string, code string) (err error) {
	from := settings.EmailAddress
	message := "From: <" + from + ">\n" +
		"To: <" + to + ">\n" +
		"Subject: Notes Overflow login code.\n\n" +
		"Hello!\nYour login code is " + code + "\n" +
		"Ignore this message if you have not requested a login.\n\n" +
		"- Notes Overflow\nThis message is sent from an unmonitored inbox."

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, settings.EmailPassword, "smtp.gmail.com"),
		from, []string{to}, []byte(message))

	return err
}
