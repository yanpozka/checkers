package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats"
)

const (
	writeWait      = 5 * time.Second     // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 2048                // Maximum message size allowed from peer.
)

var upgrader = websocket.Upgrader{} // default ReadBufferSize, WriteBufferSize 4K

type client struct {
	conn     *websocket.Conn // The websocket connection.
	send     chan []byte     // channel of outbound messages.
	done     chan struct{}
	gameID   string
	playerID string
}

func (c *client) read() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error Reading: %v\n", err)
			}
			break
		}

		c.send <- message
	}
}

func (c *client) startListenFromNats() error {
	subj := "subject-" + c.gameID
	ch := make(chan *nats.Msg, 1)

	sub, err := natsClient.ChanSubscribe(subj, ch)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if err := sub.Unsubscribe(); err != nil {
				log.Printf("Very unexpected error on Unsubscribing: %v\n", err)
			}
			c.conn.Close()
		}()

		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				log.Printf("Message from nats: %q\n", string(m.Data))
			case <-c.done:
				return
			}
		}
	}()

	return nil
}

func (c *client) write() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message := <-c.send:
			log.Printf("Resending message: %q\n", string(message))

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				log.Println("Error writing Message", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Error writing Ping Message", err)
				return
			}
		case <-c.done:
			return
		}
	}
}
