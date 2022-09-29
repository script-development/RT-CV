package models

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

	URL        string   `json:"url"`
	Method     string   `json:"method" description:"the method to use when calling the url (GET, POST, PUT, PATCH, DELETE)"`
	AddHeaders []Header `json:"addHeaders" bson:"addHeaders"`
}

// CollectionName returns the collection name of the Profile
func (*OnMatchHook) CollectionName() string {
	return "onMatchHooks"
}

// GetOnMatchHooks returns all the onMatchHooks
func GetOnMatchHooks(dbConn db.Connection, expectAtLeastOne bool) ([]OnMatchHook, error) {
	hooks := []OnMatchHook{}
	err := dbConn.Find(&OnMatchHook{}, &hooks, nil)
	if err != nil {
		return nil, err
	}

	if expectAtLeastOne && len(hooks) == 0 {
		return nil, errors.New("no on match hooks configured")
	}

	return hooks, err
}

// DataKind deinfes the kind of data that is being sent to the hook
type DataKind uint8

const (
	// DataKindMatch is the data kind for when a match is found
	DataKindMatch DataKind = iota
	// DataKindList is the data kind for when a list of cvs is matched
	DataKindList
)

func (k DataKind) contentTypeAndDataKind() (contentType string, dataKind string) {
	switch k {
	case DataKindMatch:
		contentType, dataKind = "application/json", "match"
	case DataKindList:
		contentType, dataKind = "application/json", "list"
	}
	return contentType, dataKind
}

// CallAndLogResult calls the hook defined in OnMatchHook and logs the result
func (h *OnMatchHook) CallAndLogResult(body io.Reader, dataKind DataKind, logger *log.Entry) {
	_, err := h.Call(body, dataKind)

	loggerWithFields := logger.WithField("hook", h.URL).WithField("hook_id", h.ID.Hex())
	if err != nil {
		loggerWithFields.WithError(err).Error("calling hook failed")
	} else {
		loggerWithFields.Info("hook called")
	}
}

// Call calls the hook defined in OnMatchHook
func (h *OnMatchHook) Call(body io.Reader, dataKind DataKind) (http.Header, error) {
	req, err := http.NewRequest(h.Method, h.URL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "RT-CV")

	contentTypeHeader, dataKindHeader := dataKind.contentTypeAndDataKind()
	req.Header.Set("Content-Type", contentTypeHeader)
	req.Header.Set("Data-Kind", dataKindHeader)

	for _, header := range h.AddHeaders {
		for _, value := range header.Value {
			req.Header.Add(header.Key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return req.Header, err
	}

	if resp.StatusCode >= 400 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return req.Header, fmt.Errorf("hook returned status code \"%s\" with a unreadable message, error: %s", resp.Status, err.Error())
		}

		return req.Header, fmt.Errorf("hook returned status code \"%s\" with message: %s", resp.Status, string(respBody))
	}

	return req.Header, nil
}
