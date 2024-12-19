package websocket

import (
	"bytes"
	"encoding/json"
	"log/slog"
	db "tacticstoe/internal/database"
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

	user *db.User

	queue *Queue
	hub   *Hub
	game  *Game

	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		slog.Info("Closing readPump, id: " + c.user.ID)
		c.hub.unregister <- c
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
			}
			break
		}

		r := bytes.NewReader(message)

		var clientMessage ClientMessage
		if err := json.NewDecoder(r).Decode(&clientMessage); err != nil {
			slog.Error("Error decoding message: " + err.Error())
			continue
		}

		switch clientMessage.Type {
		case "action":
			var clientActionsData ClientActionsData
			if err := mapstructure.Decode(clientMessage.Data, &clientActionsData); err != nil {
				slog.Error("Error decoding action data: " + err.Error())
				continue
			}

			if c.game == nil {
				slog.Error("Client is not part of any game")
				continue
			}

			c.game.makeAction <- GameAction{
				player: c,
				x:      clientActionsData.X,
				y:      clientActionsData.Y,
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		slog.Info("Closing writePump, id: " + c.user.ID)

		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error("Error getting writer: " + err.Error())
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				slog.Error("Error closing writer: " + err.Error())
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error("Error sending ping: " + err.Error())
				return
			}
		}
	}

}
