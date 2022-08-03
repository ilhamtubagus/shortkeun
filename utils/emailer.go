package utils

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendHTMLMail(to []string, subject string, body interface{}, templatePath string, attachmentsPath []string) error {
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")
	mailUser := os.Getenv("MAIL_USR")
	mailPasswd := os.Getenv("MAIL_PASSWD")
	mailFrom := os.Getenv("MAIL_FROM")
	if smtpHost == "" || smtpPort == "" || mailUser == "" || mailPasswd == "" || mailFrom == "" {
		return errors.New("error while loading mail configuration")
	}
	smtpPortConverted, err := strconv.Atoi(smtpPort)
	if err != nil {
		return errors.New("error while loading mail configuration [port must be a number]")
	}

	//render template
	var errParsing error
	t, errParsing := template.ParseFiles(templatePath)
	if errParsing != nil {
		log.Println(errParsing)
	}
	var templateString bytes.Buffer
	if err := t.Execute(&templateString, body); err != nil {
		log.Println(err)
	}
	htmlString := templateString.String()

	//construct mail
	mail := gomail.NewMessage()
	mail.SetHeader("From", mailFrom)
	mail.SetHeader("To", to...)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", htmlString)
	for _, v := range attachmentsPath {
		mail.Embed(v)
	}
	auth := smtp.PlainAuth("", mailUser, mailPasswd, smtpHost)
	//send mail
	mailer := gomail.Dialer{Host: smtpHost, Port: smtpPortConverted, Username: mailUser, Password: mailPasswd, Auth: auth}
	if err := mailer.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}
