package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/apex/log"
)

// SetupLogHandler changes the log handler to send errors, fatals and warns to slack
func SetupLogHandler(webhookURL string) {
	messageEnv := os.Getenv("SLACK_ENVIRONMENT")
	originalHandler := log.Log.(*log.Logger).Handler.HandleLog
	var handler log.HandlerFunc = func(entry *log.Entry) error {
		handleLogEntry(entry, webhookURL, messageEnv)
		return originalHandler(entry)
	}
	log.SetHandler(handler)
}

func handleLogEntry(entry *log.Entry, webhookURL string, messageEnv string) {
	// IMPORTANT:
	// Do not send error, fatal or warn logs in this method as it will cause a n+1 loop

	switch entry.Level {
	case log.WarnLevel, log.FatalLevel, log.ErrorLevel:
		// Continue
	default:
		return
	}

	text := fmt.Sprintf("*RTCV %s %s: %s*", messageEnv, entry.Level, strings.TrimSpace(entry.Message))

	if len(entry.Fields) != 0 {
		text += "\n\nFields:"
		for key, value := range entry.Fields {
			text += fmt.Sprintf("\n%s: %+v", key, value)
		}
	}

	type m = map[string]any
	body, err := json.Marshal(m{
		"text": "",
		"blocks": []m{{
			"type": "section",
			"text": m{"type": "mrkdwn", "text": text},
		}},
	})
	if err != nil {
		log.WithError(err).Info("Unable to marshal log entry meant for slack webhook")
		return
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(body))
	if err != nil {
		log.WithError(err).Info("Unable create post request for slack webhook")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Info("Unable to post to slack webhook")
		return
	}

	if resp.StatusCode >= 400 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Info(fmt.Sprintf("slack webhook returned status code \"%s\" with a unreadable message", resp.Status))
		} else {
			log.Info(fmt.Sprintf("slack webhook returned status code \"%s\" with message: %s", resp.Status, string(respBody)))
		}
	}
}
