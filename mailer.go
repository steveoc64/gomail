package gomail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

const mimeHeader = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n";

type Mailer struct {
	Server string
	Port   int
	User   string
	Pass   string
}

func New(server string, user, pass string) *Mailer {
	return &Mailer{
		Server: server,
		Port:   587,
		User:   user,
		Pass:   pass,
	}
}

func (m *Mailer) Send(from, to, subject, msg string) error {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", m.Server, m.Port))
	if err != nil {
		return err
	}

	if err := c.Hello(m.Server); err != nil {
		return err
	}

	config := &tls.Config{
		ServerName:         m.Server,
		InsecureSkipVerify: true,
	}
	if err := c.StartTLS(config); err != nil {
		return err
	}

	if err := c.Auth(smtp.PlainAuth("", m.User, m.Pass, m.Server)); err != nil {
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	// Send the email subject and body.
	wc, err := c.Data()
	if err != nil {
		return err
	}

	sendBytes := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n%s\n", from, to, subject, mimeHeader)
	sendBytes = sendBytes + "<html><body>\n" + msg + "\n</body></html>\n"
	if _, err := wc.Write([]byte(sendBytes)); err != nil {
		return err
	}

	err = wc.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		return err
	}
	return nil
}
