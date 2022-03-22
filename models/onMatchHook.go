package models

import (
	"io"
	"net/http"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Header is a struct that contains a http header
type Header struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

// OnMatchHook can hook onto the matching process and call API calls in case of matches
type OnMatchHook struct {
	db.M  `bson:",inline"`
	KeyID primitive.ObjectID `json:"keyId" bson:"keyId"`

	URL                  string   `json:"url"`
	Method               string   `json:"method" description:"the method to use when calling the url (GET, POST, PUT, PATCH, DELETE)"`
	AddHeaders           []Header `json:"addHeaders" bson:"addHeaders"`
	StopRemainingActions bool     `json:"stopRemainingActions" bson:"stopRemainingActions" description:"If true, the hook will stop the remaining actions after this hook such as sending an email"`
}

// CollectionName returns the collection name of the Profile
func (*OnMatchHook) CollectionName() string {
	return "onMatchHooks"
}

// Call calls the hook defined in OnMatchHook
func (h *OnMatchHook) Call(body io.Reader) error {
	req, err := http.NewRequest(h.Method, h.URL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "RT-CV")

	for _, header := range h.AddHeaders {
		for _, value := range header.Value {
			req.Header.Add(header.Key, value)
		}
	}

	// We don't really care about what the response is
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).WithField("url", h.URL).WithField("id", h.ID.Hex()).Warnf("Failed calling hook")
	}

	return nil
}
