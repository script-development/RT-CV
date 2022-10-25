package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/script-development/RT-CV/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Header is a struct that contains a http header
type Header struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

// OnMatchHook can hook onto the matching process and call API calls in case of matches
type OnMatchHook struct {
	db.M     `bson:",inline"`
	KeyID    primitive.ObjectID `json:"keyId" bson:"keyId"`
	Disabled bool               `json:"disabled" bson:"disabled"`

	URL        string   `json:"url"`
	Method     string   `json:"method" description:"the method to use when calling the url (GET, POST, PUT, PATCH, DELETE)"`
	AddHeaders []Header `json:"addHeaders" bson:"addHeaders"`
}

// CollectionName returns the collection name of the Profile
func (*OnMatchHook) CollectionName() string {
	return "onMatchHooks"
}

// GetOnMatchHooksProps contains the properties for GetOnMatchHooks
type GetOnMatchHooksProps struct {
	AllowDisabled    bool
	ExpectAtLeastOne bool
}

// GetOnMatchHooks returns all the onMatchHooks
func GetOnMatchHooks(dbConn db.Connection, props GetOnMatchHooksProps) ([]OnMatchHook, error) {
	query := bson.M{}
	if !props.AllowDisabled {
		// We do not allow disabled
		//
		// As disabled was later added to OnMatchHook it might be missing from the database
		// Because of this we check for the opposite of disabled hence this bloated query
		query["disabled"] = bson.M{"$not": bson.M{"$eq": true}}
	}

	hooks := []OnMatchHook{}
	err := dbConn.Find(&OnMatchHook{}, &hooks, query)
	if err != nil {
		return nil, err
	}

	if props.ExpectAtLeastOne && len(hooks) == 0 {
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
	_, err := h.CallWithRetry(body, dataKind)

	loggerWithFields := logger.WithField("hook", h.URL).WithField("hook_id", h.ID.Hex())
	if err != nil {
		loggerWithFields.WithError(err).Error("calling hook failed")
	} else {
		loggerWithFields.Info("hook called")
	}
}

// CallWithRetry executes (*OnMatchHook).Call() with a retry if it failes with spesific reasons
func (h *OnMatchHook) CallWithRetry(body io.Reader, dataKind DataKind) (http.Header, error) {
	reqID := primitive.NewObjectID().String()
	// do 5 retries
	var headers http.Header
	var err error
	for i := 0; i < 5; i++ {
		// Is retry, do a backoff
		switch i {
		case 1:
			time.Sleep(time.Millisecond * 100)
		case 2:
			time.Sleep(time.Millisecond * 500)
		case 3:
			time.Sleep(time.Second)
		case 4:
			time.Sleep(time.Second * 2)
		}

		if body == nil {
			body = bytes.NewReader(nil)
		}
		headers, err = h.Call(body, dataKind, reqID)
		if err == nil {
			break
		}
		statusCodeErr, ok := err.(*StatusCodeError)
		if !ok {
			break
		}

		retry := false
		switch statusCodeErr.code {
		case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			retry = true
		}

		if !retry {
			break
		}
	}
	return headers, err
}

// StatusCodeError is an error thrown by (*OnMatchHook).Call() when the status code is >= 400
type StatusCodeError struct {
	status string
	code   int
	body   []byte
}

func (e *StatusCodeError) Error() string {
	if len(e.body) == 0 {
		return fmt.Sprintf("hook returned status code \"%s\" with a unreadable message", e.status)
	}
	return fmt.Sprintf("hook returned status code \"%s\" with message: %s", e.status, string(e.body))
}

// Call calls the hook defined in OnMatchHook
func (h *OnMatchHook) Call(body io.Reader, dataKind DataKind, reqID string) (http.Header, error) {
	req, err := http.NewRequest(h.Method, h.URL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "RT-CV")
	req.Header.Set("X-Request-ID", reqID)

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
		respBody, _ := ioutil.ReadAll(resp.Body)
		return req.Header, &StatusCodeError{
			status: resp.Status,
			code:   resp.StatusCode,
			body:   respBody,
		}
	}

	return req.Header, nil
}
