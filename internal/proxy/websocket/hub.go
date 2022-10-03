package websocket

import (
	"sync"

	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
)

type hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan *domain.LogEntry

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var lock = &sync.Mutex{}

var singleInstance *hub

func GetHubInstance() *hub {
	lock.Lock()
	defer lock.Unlock()

	if singleInstance == nil {
		if singleInstance == nil {
			singleInstance = &hub{
				Broadcast:  make(chan *domain.LogEntry),
				register:   make(chan *Client),
				unregister: make(chan *Client),
				clients:    make(map[*Client]bool),
			}
			go singleInstance.run()
		}
	}

	return singleInstance
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.clientLogEntryChan)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				if client.shouldWrite(message) {
					select {
					case client.clientLogEntryChan <- message:
					default:
						close(client.clientLogEntryChan)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
