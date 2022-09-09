package controller

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RouteScraperListCVsReq contains the request data of routeScraperListCVs
type RouteScraperListCVsReq struct {
	CVs []models.CV `json:"cvs"`
}

// RouteScraperListCVsResp is the response of routeScraperListCVs
type RouteScraperListCVsResp struct{}

// CVListsHookData contains the data that is passed to the hook
type CVListsHookData struct {
	CVs              map[string]models.CV            `json:"cvs" description:"The CVs that where matched with a profile. The key is the reference number of the cv"`
	ProfilesMatchCVs map[primitive.ObjectID][]string `json:"profilesMatchCVs" description:"a map where the key is the profile ID and the value is a list of CV reference numbers it matched with"`
	KeyID            primitive.ObjectID              `json:"keyId" description:"The ID of the API key that was used to upload this CV"`
	KeyName          string                          `json:"keyName" description:"The Name of the API key that was used to upload this CV"`
	IsTest           bool                            `json:"isTest" description:"True if this hook call was manually triggered"`
}

var routeScraperListCVs = routeBuilder.R{
	Description: "Main route to scrape the CV",
	Res:         RouteScraperListCVsResp{},
	Body:        RouteScraperListCVsReq{},
	Fn: func(c *fiber.Ctx) error {
		body := RouteScraperListCVsReq{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		if len(body.CVs) == 0 {
			return errors.New("no CVs where provided in body")
		}

		reqCtx := ctx.Get(c)

		profilesCache, err := reqCtx.GetOrGenMatcherProfilesCache()
		if err != nil {
			return err
		}
		if len(profilesCache.ListProfiles) == 0 {
			// No list profiles to use for matching, so we can just return
			return c.JSON(RouteScraperListCVsResp{})
		}

		hookData := CVListsHookData{
			CVs:              map[string]models.CV{},
			ProfilesMatchCVs: map[primitive.ObjectID][]string{},
		}

		for _, cv := range body.CVs {
			if cv.PersonalDetails.Zip == "" {
				continue
			}

			cvZip, validZip := cv.PersonalDetails.ZipAsNr()
			if !validZip {
				continue
			}
			cvRef := cv.ReferenceNumber

			cvMatch := false
			for _, profile := range profilesCache.ListProfiles {
				for _, zipCode := range profile.Zipcodes {
					if zipCode.IsWithinCithAndArea(cvZip) {
						cvMatch = true
						hookData.ProfilesMatchCVs[profile.ID] = append(hookData.ProfilesMatchCVs[profile.ID], cvRef)
						break
					}
				}
			}

			if cvMatch {
				hookData.CVs[cvRef] = cv
			}
		}

		if len(hookData.CVs) == 0 {
			return c.JSON(RouteScraperListCVsResp{})
		}

		dataForHook, err := json.Marshal(hookData)
		if err != nil {
			return err
		}

		go func(dataForHook []byte) {
			hooks, err := models.GetOnMatchHooks(reqCtx.DBConn, true)
			if err != nil {
				reqCtx.Logger.WithError(err).Error("Finding on match hooks failed")
				return
			}
			if len(hooks) == 0 {
				reqCtx.Logger.Info("No hooks configured to send matched list CVs to")
				return
			}

			for _, hook := range hooks {
				hook.CallAndLogResult(bytes.NewReader(dataForHook), models.DataKindList, reqCtx.Logger)
			}
		}(dataForHook)

		return c.JSON(RouteScraperListCVsResp{})
	},
}
