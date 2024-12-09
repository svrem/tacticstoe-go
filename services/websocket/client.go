package websocket_service

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type ClientMessage struct {
	Type string            `json:"type"`
	Data ClientMessageData `json:"data"`
}
type ClientMessageData map[string]interface{}

type ClientActionsData struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Client struct {
	conn *websocket.Conn

	elo_rating int

	id string

	queue *Queue

	game *Game

	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		if c.game != nil {
			slog.Info("Requesting unregistering client from the game.")
			c.game.unregister <- c
		}
		if c.queue != nil {
			slog.Info("Unregistering client from the queue.")
			c.queue.unregister <- c
		}
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Unexpected close error: " + err.Error())
			} else {
				slog.Info("Connection closed: " + err.Error())
			}
			break
		}

		r := bytes.NewReader(message)

		var client_message ClientMessage
		if err := json.NewDecoder(r).Decode(&client_message); err != nil {
			slog.Error("Error decoding message: " + err.Error())
			continue
		}

		switch client_message.Type {
		case "action":
			var client_actions_data ClientActionsData
			if err := mapstructure.Decode(client_message.Data, &client_actions_data); err != nil {
				slog.Error("Error decoding action data: " + err.Error())
				continue
			}

			if c.game == nil {
				slog.Error("Client is not part of any game")
				continue
			}

			c.game.make_action <- GameAction{
				player: c,
				x:      client_actions_data.X,
				y:      client_actions_data.Y,
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			slog.Info("Sending message: " + string(message))
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error("Error getting writer: " + err.Error())
				return
			}
			w.Write(message)

			// // Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				slog.Error("Error closing writer: " + err.Error())
				return
			}
		case <-ticker.C:
			slog.Info("Sending ping")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Error sending ping: " + err.Error())
				return
			}
		}
	}

}
