package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/script-development/RT-CV/helpers/routeBuilder"
)

var routeControlEventsWS = routeBuilder.R{
	Description: "Used by the dashboard to get continues updates from the backend\n\n" +
		"Note this route returns a websocket connection and not an empty object",
	Res: struct{}{},
	Fn: func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			// This is not request for a websocket connection
			return fiber.ErrUpgradeRequired
		}

		return routeControlEventsWSHandler(c)
	},
}

var routeControlEventsWSHandler = websocket.New(func(c *websocket.Conn) {
	// Somtimes due to a network setup a websocket might be automatically closed
	// If there are no messages send / received so we send a ping message every 30 seconds
	keepAliveTicker := time.NewTicker(time.Second * 30)
	for {
		<-keepAliveTicker.C
		err := c.WriteMessage(websocket.PingMessage, []byte("PING"))
		if err != nil {
			break
		}
	}
	keepAliveTicker.Stop()

	// // websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	// var (
	// 	mt  int
	// 	msg []byte
	//  err error
	// )

	// for {
	// 	mt, msg, err = c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", msg)

	// 	err = c.WriteMessage(mt, msg)
	// 	if err != nil {
	// 		log.Println("write:", err)
	// 		break
	// 	}
	// }
})
