package services

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func SendMail(to, message, title string) {
	from := "garifullin.ernur@mail.ru"
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)

	m.SetBody("text/plain", message)

	d := gomail.NewDialer(
		"smtp.mail.ru", 465,
		"garifullin.ernur@mail.ru",
		"S666jY9t16E9TUh1RUgK",
	)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}
