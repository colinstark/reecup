package server

import (
	"encoding/json"
	"log"
	"net/http"

	"reecup/game"
	"sync"

	"github.com/gorilla/websocket"
)

type userID = string
type gameID = string

type GameServer struct {
	upgrader    websocket.Upgrader
	Connections map[userID]*websocket.Conn
	Users       map[userID]User      `json:"users"`
	Games       map[gameID]game.Game `json:"games"`
	userMutex   sync.RWMutex
	gameMutex   sync.RWMutex
}

type User struct {
	Name string `json:"name"`
}

func NewGameServer() *GameServer {
	return &GameServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Connections: make(map[userID]*websocket.Conn),
		Users:       make(map[userID]User),
		Games:       make(map[gameID]game.Game),
	}
}

func (s *GameServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("NEW CONNECTION")

	userID := getUserIDFromRequest(r)

	if userID == "" {
		log.Println("No user ID provided")
		newUserID := generateUserID()
		log.Println("Sending ID", newUserID)
		conn.WriteJSON(map[string]any{
			"type":    "user_id_commissioned",
			"user_id": newUserID,
		})
	} else {
		s.registerClient(userID, conn)
	}

	// Main loop
	for {
		_, message, err := conn.ReadMessage()
		log.Println("GOT MESSAGE")
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var data map[string]any
		if err := json.Unmarshal(message, &data); err != nil {
			log.Println("JSON error:", err)
			continue
		}

		switch data["instruction"] {
		case "login":
			s.handleLogin(data, userID)
		case "logout":
			s.handleLogout(data, userID)
		case "update_name":
			s.handleUpdateName(data, userID)
		case "new_game":
			s.handleNewGame(userID)
		case "list_games":
			s.handleListGames(data, userID)
		case "start_game":
			s.handleStartGame(data, userID)
		case "cancel_game":
			s.handleCancelGame(data, userID)
		case "join_game":
			s.handleJoinGame(data, userID)
		case "draw_stone":
			s.handleDrawStone(data, userID)
		case "finish_turn":
			s.handleFinishTurn(data, userID)
		case "get_deck":
			s.handleGetDeck(data, userID)
		case "update_cursor":
			// s.handleUpdateCursor(data, userID)
		default:
			s.handleError("general", "Invalid Event", userID)
		}
	}

	s.unregisterClient(userID, conn)
}

func (s *GameServer) registerClient(userID string, conn *websocket.Conn) {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	s.Connections[userID] = conn

	if _, exists := s.Users[userID]; !exists {
		s.Users[userID] = User{
			Name: "", // will be set later async
		}
	}
}

func (s *GameServer) unregisterClient(userID string, conn *websocket.Conn) {
	_ = conn
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	delete(s.Connections, userID)
}

func (s *GameServer) handleError(instruction string, error string, targetID string) {
	conn := s.Connections[targetID]
	if conn != nil {
		conn.WriteJSON(map[string]any{
			"instruction": instruction,
			"error":       error,
		})
	}
}

func (s *GameServer) Send(target string, instruction string, data map[string]any) {
	conn := s.Connections[target]
	if conn != nil {
		json, err := marshal(instruction, data)
		if err == nil {
			conn.WriteMessage(websocket.TextMessage, json)
		}
	}
}

func (s *GameServer) Broadcast(instruction string, data map[string]any) {
	json, err := marshal(instruction, data)
	if err != nil {
		return
	}

	for _, conn := range s.Connections {
		conn.WriteMessage(websocket.TextMessage, json)
	}
}

func (s *GameServer) BroadcastInGame(gameID string, instruction string, data map[string]any) {
	s.gameMutex.RLock()
	game, exists := s.Games[gameID]
	if !exists {
		s.gameMutex.RUnlock()
		return
	}
	s.gameMutex.RUnlock()

	json, err := marshal(instruction, data)
	if err != nil {
		return
	}

	// Send message only to players in this game
	for _, player := range game.Players {
		if conn, exists := s.Connections[player.ID]; exists {
			conn.WriteMessage(websocket.TextMessage, json)
		}
	}
}

func marshal(instruction string, data map[string]any) ([]byte, error) {
	data["instruction"] = instruction
	responseJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Marshal error:", err)
		return []byte{}, err
	}

	return responseJSON, nil
}

func getUserIDFromRequest(r *http.Request) string {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Println("No user ID provided")
	}

	return userID
}
