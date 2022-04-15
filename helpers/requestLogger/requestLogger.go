package requestLogger

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/script-development/RT-CV/controller/ctx"
)

// New creates a new fiber logger middleware
func New() fiber.Handler {
	var timestamp atomic.Value
	timestamp.Store(time.Now().Format("15:04:05"))

	go func() {
		for {
			time.Sleep(time.Second)
			timestamp.Store(time.Now().Format("15:04:05"))
		}
	}()

	return func(c *fiber.Ctx) (err error) {
		path := c.Path()
		if path == "/api/v1/health" {
			// This URL is request a lot it and there isn't much value logging it
			// If anything it bloats the logs
			return c.Next()
		}

		start := time.Now()
		err = c.Next()
		reqDuration := time.Since(start).Microseconds()

		resp := c.Response()

		// ${time} ${status} - ${latency} ${method} ${path}\nfmt.Print(timestamp)
		fmt.Print(timestamp.Load().(string))
		fmt.Print(" | ")
		fmt.Print(resp.StatusCode())
		fmt.Print(" | ")
		if reqDuration >= 2_000 {
			// When a request takes more than 2ms we log it as ms
			reqDurationStr := strconv.FormatInt(reqDuration/1_000, 10)
			fmt.Print(strings.Repeat(" ", 5-len(reqDurationStr)))
			fmt.Print(reqDurationStr)
			fmt.Print("ms | ")
		} else {
			reqDurationStr := strconv.FormatInt(reqDuration, 10)
			fmt.Print(strings.Repeat(" ", 5-len(reqDurationStr)))
			fmt.Print(reqDurationStr)
			fmt.Print("µs | ")
		}

		apikey := ctx.Get(c).Key
		if apikey != nil {
			nameLen := len(apikey.Name)
			if nameLen > 15 {
				fmt.Print(apikey.Name[:15])
			} else {
				fmt.Print(strings.Repeat(" ", 15-nameLen))
				fmt.Print(apikey.Name)
			}
		} else {
			ip := c.IP()
			fmt.Print(strings.Repeat(" ", 15-len(ip)))
			fmt.Print(ip)
		}
		fmt.Print(" | ")
		method := c.Method()
		fmt.Print(method)
		fmt.Print(strings.Repeat(" ", 7-len(method)))
		fmt.Print("| ")
		fmt.Print(path)
		if err != nil {
			fmt.Print(" | ERR: ")
			fmt.Print(err)
		}
		fmt.Print("\n")

		return err
	}
}
