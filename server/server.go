package server

import (
	"encoding/json"
	"log"
	"net/http"
	"reecup/game"
	"sync"

	"github.com/gorilla/websocket"
)

type GameServer struct {
	upgrader  websocket.Upgrader
	clients   map[string]*websocket.Conn
	users     map[string]User
	games     map[string]game.Game
	userMutex sync.RWMutex
	gameMutex sync.RWMutex
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewGameServer() *GameServer {
	return &GameServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[string]*websocket.Conn),
		users:   make(map[string]User),
		games:   make(map[string]game.Game),
	}
}

func (s *GameServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	userID := getUserIDFromRequest(r)
	if userID == "" {
		log.Println("No user ID provided")
		return
	}

	s.registerClient(userID, conn)

	// Main loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("JSON error:", err)
			continue
		}

		switch data["instruction"] {
		case "login":
		// s.handleLogin(data, userID)
		case "new_game":
		// s.handleNewGame(data, userID)
		case "list_games":
		// s.handleListGames(data, userID)
		case "join_game":
		// s.handleJoinGame(data, userID)
		case "draw_stone":
		// s.handleDrawStone(data, userID)
		case "finish_turn":
			// s.handleFinishTurn(data, userID)
		case "update_cursor":
			s.handleUpdateCursor(data, userID)
		case "get_deck":
			s.handleGetDeck(data, userID)
		default:
			s.handleError("general", "Invalid Event", userID)
		}
	}

	s.unregisterClient(userID, conn)
}

func (s *GameServer) registerClient(userID string, conn *websocket.Conn) {
	s.userMutex.Lock()
	s.clients[userID] = conn
	s.userMutex.Unlock()
}

func (s *GameServer) unregisterClient(userID string, conn *websocket.Conn) {
	_ = conn
	s.userMutex.Lock()
	delete(s.clients, userID)
	s.userMutex.Unlock()
}

func (s *GameServer) handleError(instruction string, error string, senderID string) {
	responseJSON, _ := json.Marshal(map[string]interface{}{
		"instruction": instruction,
		"error":       error,
	})

	s.clients[senderID].WriteMessage(websocket.TextMessage, responseJSON)
}

func getUserIDFromRequest(r *http.Request) string {
	return r.URL.Query().Get("userID")
}
