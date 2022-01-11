package controller

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
	"github.com/script-development/RT-CV/models"
)

// RouteGetExampleExampleAttachmentPDFBody is the body data for the route below
type RouteGetExampleExampleAttachmentPDFBody struct {
	Options *models.PdfOptions `json:"options"`
}

var routeGetExampleAttachmentPDF = routeBuilder.R{
	Description: "Download an example email attachment PDF, with optional options",
	Body:        RouteGetExampleExampleAttachmentPDFBody{},
	Res:         []byte{},
	Fn: func(c *fiber.Ctx) error {
		body := RouteGetExampleExampleAttachmentPDFBody{}
		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		cv := models.ExampleCV()
		pdfFile, err := cv.GetPDF(body.Options, nil)
		if err != nil {
			return err
		}

		defer func() {
			pdfFile.Close()
			os.Remove(pdfFile.Name())
		}()

		pdfFile.Seek(0, 0)
		pdfBytes, err := ioutil.ReadAll(pdfFile)
		if err != nil {
			return errors.New("unable to read the generated pdf file")
		}

		c.Set("Content-type", "application/pdf")
		return c.Send(pdfBytes)
	},
}
