package user

import (
	"books/app"
	"books/app/common"
	"gopkg.in/gomail.v2"
)

func SendEmailWithResetLink(user UserModel) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", app.TEMP_EMAIL, "Books Team")
	m.SetAddressHeader("To", user.Email, user.Username)
	m.SetHeader("Subject", "Reset Password")

	m.SetBody("text/plain", "Token for request: "+common.GenToken(user.ID))

	d := gomail.NewPlainDialer("smtp.yandex.ru", 587, app.TEMP_EMAIL, app.TEMP_EMAIL_PASSWORD)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
