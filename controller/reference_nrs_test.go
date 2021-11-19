package controller

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/script-development/RT-CV/db"
	"github.com/script-development/RT-CV/helpers/jsonHelpers"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/mock"
	"github.com/script-development/RT-CV/models"
	. "github.com/stretchr/testify/assert"
)

func TestScannedReferenceNrs(t *testing.T) {
	router := newTestingRouter(t)

	newParsedCvRef := func(nr string, insertionTime time.Time) *models.ParsedCVReference {
		return &models.ParsedCVReference{
			M:               db.NewM(),
			ReferenceNumber: nr,
			InsertionDate:   jsonHelpers.RFC3339Nano(insertionTime),
			KeyID:           mock.Key1.ID,
		}
	}

	err := router.db.UnsafeInsert(
		newParsedCvRef("1", time.Now().Add(-(time.Minute*30))),
		newParsedCvRef("2", time.Now().Add(-(time.Minute*90))),
		newParsedCvRef("3", time.Now().AddDate(0, 0, -2)),
		newParsedCvRef("4", time.Now().AddDate(0, 0, -4)),
		newParsedCvRef("5", time.Now().AddDate(0, 0, -8)),  // 1 week + 1 day
		newParsedCvRef("6", time.Now().AddDate(0, 0, -15)), // 2 weeks + 1 day
		newParsedCvRef("7", time.Now().AddDate(0, 0, -21)), // 3 weeks + 1 day
	)
	NoError(t, err)

	doTest := func(t *testing.T, routeSuffix string, expectRefNrs ...string) {
		res, body := router.MakeRequest(routeBuilder.Get, "/api/v1/scraper/scannedReferenceNrs"+routeSuffix, TestReqOpts{})
		Equal(t, 200, res.StatusCode, string(body))

		refNrs := []string{}
		err = json.Unmarshal(body, &refNrs)
		NoError(t, err)

		Equal(t, expectRefNrs, refNrs)
	}

	t.Run("all reference nrs", func(t *testing.T) {
		doTest(t, "", "1", "2", "3", "4", "5", "6", "7")
	})

	t.Run("reference nrs since hours", func(t *testing.T) {
		doTest(t, "/since/hours/1", "1")
		doTest(t, "/since/hours/2", "1", "2")
	})

	t.Run("reference nrs since days", func(t *testing.T) {
		doTest(t, "/since/days/1", "1", "2")
		doTest(t, "/since/days/3", "1", "2", "3")
	})

	t.Run("reference nrs since weeks", func(t *testing.T) {
		doTest(t, "/since/weeks/1", "1", "2", "3", "4")
		doTest(t, "/since/weeks/3", "1", "2", "3", "4", "5", "6")
	})
}
