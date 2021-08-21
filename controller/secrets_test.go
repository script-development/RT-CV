package controller

import (
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestSecretRoutes(t *testing.T) {
	app := newTestingRouter()

	contents := `{"key":"value"}`
	valueKey := "test1"
	encryptionKey := "very-secret-key-of-minimal-16-chars"

	// Insert key works
	route := fmt.Sprintf("/v1/scraper/secret/%v/%v", valueKey, encryptionKey)
	res := app.MakeRequest(t, Post, route, TestReqOpts{
		Body: []byte(contents),
	})
	body, err := ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, contents, string(body))

	// Get key works
	res = app.MakeRequest(t, Get, route, TestReqOpts{
		Body: []byte(contents),
	})
	body, err = ioutil.ReadAll(res.Body)
	NoError(t, err)
	Equal(t, contents, string(body))
}
