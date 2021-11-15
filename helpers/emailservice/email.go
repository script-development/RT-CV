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

// ErrNoConf = email service not configured
var ErrNoConf = errors.New("email service not configured")

// Setup sets up the email sender
func Setup(conf EmailServerConfiguration, onMailSend func(err error)) error {
	if conf.Host == "" || conf.From == "" {
		log.Warn("Email not configured (EMAIL_HOST and EMAIL_FROM must be set), DISABELING EMAIL SUPPORT")
		go func() {
			for data := range ch {
				log.Infof("sending no mail to %v as email server is not configured", data.To)
				if onMailSend != nil {
					onMailSend(ErrNoConf)
				}
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
				retryCount := 0
				for retryCount < 4 {
					if retryCount > 0 {
						log.Infof("retrying sending mail to %v", e.To)
					} else {
						log.Infof("sending mail to %s", e.To)
					}

					e.From = from
					err := p.Send(e, 10*time.Second)
					if onMailSend != nil {
						onMailSend(err)
					}

					if err == nil {
						break
					}
					log.WithError(err).Error("sending email")
					retryCount++
				}
			}
		}(conf.From)
	}

	log.Info("Email service running..")
	return nil
}
