package websockethandlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

type ClientsStorage map[int64]map[*Client]bool

func (cs ClientsStorage) add(lotID int64, client *Client) {
	mm, ok := cs[lotID]
	if !ok {
		mm = make(map[*Client]bool)
		cs[lotID] = mm
	}
	mm[client] = true
}

func (cs ClientsStorage) delete(lotID int64, client *Client) {
	delete(cs[lotID], client)
	if len(cs[lotID]) == 0 {
		delete(cs, lotID)
	}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	lotID int64
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan BroadcastWithID
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		message, ok := <-c.send
		_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if !ok {
			// The hub closed the channel.
			_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, message.Data); err != nil {
			close(c.send)
			c.hub.clients.delete(c.lotID, c)
			fmt.Printf("removed client from lot id: %d , total clients %d\n", c.lotID, len(c.hub.clients[c.lotID]))
			return
		}

	}
}
