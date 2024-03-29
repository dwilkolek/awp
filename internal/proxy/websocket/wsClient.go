package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10 * 1024 * 1024 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	clientLogEntryChan chan *domain.LogEntry

	serviceSubscriptions sync.Map
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
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
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if strings.HasPrefix(string(message), "sub:") {
			c.serviceSubscriptions.Store(strings.Replace(string(message), "sub:", "", -1), true)
		}
		if strings.HasPrefix(string(message), "del:") {
			c.serviceSubscriptions.Store(strings.Replace(string(message), "del:", "", -1), true)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.clientLogEntryChan:
			// if c.shouldWrite(message) {
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(toBytes(message))

			// Add queued chat messages to the current websocket message.
			n := len(c.clientLogEntryChan)
			for i := 0; i < n; i++ {
				next := <-c.clientLogEntryChan
				// if c.shouldWrite(next) {
				w.Write(newline)
				w.Write(toBytes(next))
				// }
			}

			if err := w.Close(); err != nil {
				return
			}
			// }
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			fmt.Println("ticker")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) shouldWrite(entry *domain.LogEntry) bool {
	value, _ := client.serviceSubscriptions.LoadOrStore(entry.Service, false)
	if entry.Service == "" {
		log.Panicf("Entry without Service set - shouldn't happen %+v\n", entry)
	}
	if bValue, ok := value.(bool); ok {
		return bValue
	}
	return false
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: GetHubInstance(), conn: conn, clientLogEntryChan: make(chan *domain.LogEntry), serviceSubscriptions: sync.Map{}}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func toBytes(entry *domain.LogEntry) []byte {
	b, err := json.Marshal(entry)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}

	return b

}
