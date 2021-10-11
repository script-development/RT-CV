package emailservice

import (
	"errors"
	"net/smtp"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/jordan-wright/email"
)

var ch = make(chan *email.Email)

// SendMail sends an email based on the given content
func SendMail(content *email.Email) {
	if content == nil {
		return
	}
	ch <- content
}

// Setup sets up the email sender
func Setup(identity, username, password, host, port, from string) error {
	if host == "" || from == "" {
		log.Warn("Email not configured (EMAIL_HOST and EMAIL_FROM must be set), DISABELING EMAIL SUPPORT")
		go func() {
			for {
				<-ch
			}
		}()
		return nil
	}
	if port == "" {
		port = "25"
	} else {
		parsedPort, err := strconv.Atoi(port)
		if err != nil || parsedPort <= 0 {
			return errors.New("invalid port number " + port)
		}
	}

	poolSize := 4

	p, err := email.NewPool(
		host+":"+port,
		poolSize,
		smtp.PlainAuth(identity, username, password, host),
	)
	if err != nil {
		return err
	}

	for i := 0; i < poolSize; i++ {
		go func(from string) {
			for e := range ch {
				e.From = from
				err := p.Send(e, 10*time.Second)
				if err != nil {
					log.WithError(err).Error("Error sending email")
				}
			}
		}(from)
	}

	log.Info("Email service running..")
	return nil
}
