package controller

import (
	"encoding/json"
	"sync"
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
	closer := make(chan struct{})

	c.SetCloseHandler(func(int, string) error {
		close(closer)
		return nil
	})

	// Somtimes due to a network setup a websocket might be automatically closed
	// If there are no messages send / received so we send a ping message every 30 seconds
	keepAliveTicker := time.NewTicker(time.Second * 30)
	go func() {
		for range keepAliveTicker.C {
			err := c.WriteMessage(websocket.PingMessage, []byte("PING"))
			if err != nil {
				close(closer)
			}
		}
	}()

	removeListenerFn := dashboardListeners.addListener(func(jsonData []byte) {
		err := c.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			close(closer)
		}
	})

	<-closer

	keepAliveTicker.Stop()
	removeListenerFn()
})

// events contains all listeners
var dashboardListeners eventListeners

// eventListener contains a single event listener
type eventListener func(jsonData []byte)

// eventListeners is a list of event listeners
type eventListeners struct {
	m sync.Mutex

	// A pointer to eventListener is used so we can compare array entries
	// See addListener for more info
	listeners []*eventListener
}

func (e *eventListeners) addListener(listener eventListener) func() {
	e.m.Lock()
	defer e.m.Unlock()
	if listener == nil {
		return func() {}
	}

	listenerPtr := &listener
	e.listeners = append(e.listeners, listenerPtr)
	return func() {
		e.m.Lock()
		for idx, listener := range e.listeners {
			if listener == listenerPtr {
				e.listeners = append(e.listeners[:idx], e.listeners[idx+1:]...)
				break
			}
		}
		e.m.Unlock()
	}
}

func (e *eventListeners) publish(data interface{}) error {
	e.m.Lock()
	defer e.m.Unlock()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, listener := range e.listeners {
		(*listener)(bytes)
	}
	return nil
}
