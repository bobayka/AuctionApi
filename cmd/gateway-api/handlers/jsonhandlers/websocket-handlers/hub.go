package websockethandlers

import (
	"fmt"
)

const AllLotsID = -1

type BroadcastWithID struct {
	LotID int64
	Data  []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients ClientsStorage

	// Inbound messages from the clients.
	Broadcast chan BroadcastWithID

	// Register requests from the clients.
	register chan *Client
}

func newHub() *Hub {
	return &Hub{
		Broadcast: make(chan BroadcastWithID),
		register:  make(chan *Client),
		clients:   make(ClientsStorage),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients.add(client.lotID, client)
			fmt.Printf("added client on lot id: %d, total clients: %d\n", client.lotID, len(h.clients[client.lotID]))
		case message := <-h.Broadcast:
			ids := [2]int64{AllLotsID, message.LotID}
			for _, v := range ids {
				for client := range h.clients[v] {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients[client.lotID], client)
					}
				}
			}
		}
	}
}
