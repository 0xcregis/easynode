package service

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestConn(t *testing.T) {
	var addr = flag.String("addr", "localhost:9003", "http service address")
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/api/store/ws/afba013c-0072-4592-b8cd-304fa456f76e"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func Test1(t *testing.T) {

	t.Log(strings.ToLower("0xA16F524a804BEaED0d791De0aa0b5836295A2a84"))

	t.Log(strings.ToLower("0xc03fe617F10286AE63568B5c567C81ce8DF84D38"))
}
