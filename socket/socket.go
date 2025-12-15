package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowed := map[string]bool{
			"https://privmail.com":        true,
			"https://app.privmail.com":    true,
			"https://socket.privmail.com": true,
			"http://localhost:3000":       true,
		}
		if origin == "" {
			return false
		}
		return allowed[origin]
	},
	EnableCompression: true,
}

type Client struct {
	ID            string
	Email         string
	Conn          *websocket.Conn
	Ctx           context.Context
	Cancel        context.CancelFunc
	SourceId      string
	Authenticated bool
	LastActivity  time.Time
	WriteMu       sync.Mutex
}

var (
	userClients  = make(map[string][]*Client)
	allClients   = make(map[string]*Client)
	mu           sync.RWMutex
	nextClientID = 1
)

func RunWebSocketServer() {
	http.HandleFunc("/", handleWS)

	fmt.Println("WebSocket server running on :1200")
	err := http.ListenAndServe(":1200", nil)
	if err != nil {
		log.Fatalf("WebSocket server error: %v", err)
	}
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	mu.Lock()
	clientID := fmt.Sprintf("c%d", nextClientID)
	nextClientID++

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		ID:           clientID,
		Conn:         conn,
		Ctx:          ctx,
		Cancel:       cancel,
		LastActivity: time.Now(),
	}

	allClients[clientID] = client
	mu.Unlock()

	fmt.Println("Websocket client connected:", client.ID)

	go keepAliveLoop(client)

	go listen(client)
}

func assignClientEmail(c *Client, email, sourceId string) {
	mu.Lock()
	defer mu.Unlock()

	for _, cl := range userClients[email] {
		if cl.ID == c.ID {
			return
		}
	}

	c.Email = email
	c.SourceId = sourceId
	c.Authenticated = true
	userClients[email] = append(userClients[email], c)

	fmt.Printf("Client %s authenticated as %s\n", c.ID, email)
}

func removeClient(c *Client) {
	fmt.Println(`removing client from room`)
	mu.Lock()
	defer mu.Unlock()

	if c.Cancel != nil {
		c.Cancel()
	}

	if c.Email != "" {
		list := userClients[c.Email]
		for i, cl := range list {
			if cl.ID == c.ID {
				userClients[c.Email] = append(list[:i], list[i+1:]...)
				break
			}
		}
		if len(userClients[c.Email]) == 0 {
			delete(userClients, c.Email)
		}
	}

	delete(allClients, c.ID)
}

func listen(c *Client) {
	defer func() {
		fmt.Println("Client disconnected:", c.ID)
		removeClient(c)
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Printf("Message from %s: %s\n", c.ID, msg)

		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			fmt.Println(err)
			return
		}

		action := data["action"].(string)

		c.LastActivity = time.Now()

		switch action {
		case "joinRoom":
			joinRoom(c, data)
		}
	}
}

func keepAliveLoop(c *Client) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	idleLimit := 20 * time.Second

	for {
		select {
		case <-ticker.C:
			if time.Since(c.LastActivity) >= idleLimit {
				err := safeWrite(c, []byte{})
				if err != nil {
					fmt.Println("KeepAlive write failed:", err)
					c.Conn.Close()
					return
				}
				fmt.Println("Sent keepAlive to:", c.ID)
			}

		case <-c.Ctx.Done():
			return
		}
	}
}

func joinRoom(c *Client, data map[string]interface{}) {
	roomName := data["roomName"].(string)

	assignClientEmail(c, roomName, "")
}

func safeWrite(c *Client, data []byte) error {
	c.WriteMu.Lock()
	defer c.WriteMu.Unlock()
	err := c.Conn.WriteMessage(websocket.BinaryMessage, data)
	if err == nil {
		c.LastActivity = time.Now()
	}
	return err
}
