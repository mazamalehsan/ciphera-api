package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowed := map[string]bool{
			"http://localhost:3000": true,
			"https://privmail.com": true,
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
	UserId        string
	Conn          *websocket.Conn
	Ctx           context.Context
	Cancel        context.CancelFunc
	SourceId      string
	Authenticated bool
	LastActivity  time.Time
	WriteMu       sync.Mutex
}

type Group struct {
	Name    string
	Clients []string
}

var (
	userClients  = make(map[string][]*Client)
	allClients   = make(map[string]*Client)
	groups       = make(map[string]*Group)
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

func removeClient(c *Client) {
	fmt.Println(`removing client from room`)
	mu.Lock()
	defer mu.Unlock()

	if c.Cancel != nil {
		c.Cancel()
	}

	if c.UserId != "" {
		list := userClients[c.UserId]
		for i, cl := range list {
			if cl.ID == c.ID {
				userClients[c.UserId] = append(list[:i], list[i+1:]...)
				break
			}
		}
		if len(userClients[c.UserId]) == 0 {
			delete(userClients, c.UserId)
		}
	}

	delete(allClients, c.ID)
	removeClientFromRooms(c.ID)
}

func removeClientFromRooms(clientID string) {
	for roomName, grp := range groups {
		for i, id := range grp.Clients {
			if id == clientID {
				grp.Clients = append(grp.Clients[:i], grp.Clients[i+1:]...)
				break
			}
		}

		if len(grp.Clients) == 0 {
			delete(groups, roomName)
		}
	}
}

func listen(c *Client) {
	defer func() {
		fmt.Println("Client disconnected", c.ID)
		removeClient(c)
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(msg, &data); err != nil {
			fmt.Println(err)
			return
		}

		action, _ := data["action"].(string)
		c.LastActivity = time.Now()

		switch action {
		case "auth":
			handleAuth(c, data)
		case "joinRoom":
			joinRoom(c, data)
		case "sendMessage":
			handleSendMessage(c, data)
		}
	}
}

func handleAuth(c *Client, data map[string]interface{}) {
	userId, _ := data["userId"].(string)
	if userId == "" {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	c.UserId = userId
	c.Authenticated = true
	userClients[userId] = append(userClients[userId], c)

	fmt.Printf("Client %s authenticated as user %s\n", c.ID, userId)
}

func handleSendMessage(c *Client, data map[string]interface{}) {
	if !c.Authenticated {
		return
	}

	// Relay the message to recipient(s)
	// For DM: send to the "to" user's connected clients
	// For group: emit to the room

	groupId, _ := data["groupId"].(string)

	if groupId != "" {
		// Group message — emit to room
		EmitToRoom(groupId, data)
	} else {
		// DM — send to specific user
		toUserId, _ := data["to"].(string)
		if toUserId == "" {
			return
		}
		EmitToUser(toUserId, data)
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
	roomName, ok := data["roomName"].(string)
	if !ok || roomName == "" {
		return
	}

	mu.Lock()
	defer mu.Unlock()

	grp, exists := groups[roomName]
	if !exists {
		grp = &Group{
			Name:    roomName,
			Clients: []string{},
		}
		groups[roomName] = grp
	}

	// avoid duplicate join
	for _, id := range grp.Clients {
		if id == c.ID {
			return
		}
	}

	grp.Clients = append(grp.Clients, c.ID)

	fmt.Printf("Client %s joined room %s\n", c.ID, roomName)
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

func EmitToUser(userId string, payload interface{}) {
	mu.RLock()
	clients, exists := userClients[userId]
	if !exists || len(clients) == 0 {
		mu.RUnlock()
		return
	}

	clientsCopy := make([]*Client, len(clients))
	copy(clientsCopy, clients)
	mu.RUnlock()

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for _, client := range clientsCopy {
		if err := safeWrite(client, data); err != nil {
			fmt.Println("EmitToUser failed for client", client.ID, ":", err)
		}
	}
}

func EmitToRoom(roomName string, payload interface{}) {
	mu.RLock()
	group, exists := groups[roomName]
	if !exists || len(group.Clients) == 0 {
		mu.RUnlock()
		return
	}

	clientIDs := make([]string, len(group.Clients))
	copy(clientIDs, group.Clients)
	mu.RUnlock()

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for _, clientID := range clientIDs {
		mu.RLock()
		client := allClients[clientID]
		mu.RUnlock()

		if client == nil {
			continue
		}

		if err := safeWrite(client, data); err != nil {
			fmt.Println("EmitToRoom failed for client", clientID, ":", err)
		}
	}
}
