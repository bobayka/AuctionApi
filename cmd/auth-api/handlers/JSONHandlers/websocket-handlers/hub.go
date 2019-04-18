package websocket_handlers

import "fmt"

type broadcastWithID struct {
	lotID int64
	data  []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients ClientsStorage

	// Inbound messages from the clients.
	broadcast chan broadcastWithID

	// Register requests from the clients.
	register chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast: make(chan broadcastWithID),
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
		case message := <-h.broadcast:
			for client := range h.clients[message.lotID] {
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
