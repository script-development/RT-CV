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

// EmailServerConfiguration contains the configuration for the email server
type EmailServerConfiguration struct {
	// Hostname and Port of the smtp server
	Host string
	Port string

	// Authentication fields
	Identity string
	Username string
	Password string

	// The email address to send emails from
	From string
}

// Setup sets up the email sender
func Setup(conf EmailServerConfiguration, onMailSend func(err error)) error {
	if conf.Host == "" || conf.From == "" {
		log.Warn("Email not configured (EMAIL_HOST and EMAIL_FROM must be set), DISABELING EMAIL SUPPORT")
		go func() {
			for {
				<-ch
			}
		}()
		return nil
	}
	if conf.Port == "" {
		conf.Port = "25"
	} else {
		parsedPort, err := strconv.Atoi(conf.Port)
		if err != nil || parsedPort <= 0 {
			return errors.New("invalid port number " + conf.Port)
		}
	}

	poolSize := 4

	p, err := email.NewPool(
		conf.Host+":"+conf.Port,
		poolSize,
		smtp.PlainAuth(conf.Identity, conf.Username, conf.Password, conf.Host),
	)
	if err != nil {
		return err
	}

	for i := 0; i < poolSize; i++ {
		go func(from string) {
			for e := range ch {
				e.From = from
				err := p.Send(e, 10*time.Second)
				onMailSend(err)
				if err != nil {
					log.WithError(err).Error("Error sending email")
				}
			}
		}(conf.From)
	}

	log.Info("Email service running..")
	return nil
}
