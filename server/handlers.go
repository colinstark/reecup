package server

import (
	"encoding/json"
	"log"
	"reecup/game"

	"github.com/gorilla/websocket"
)

func (s *GameServer) handleLogin(data map[string]interface{}) {
}

func (s *GameServer) handleLogin(data map[string]interface{}) {
}


func (s *GameServer) handleGetDeck(data map[string]interface{}, userID string) {
	deck := game.CreateDeck()

	responseJSON, _ := json.Marshal(map[string]interface{}{
		"instruction": "new_deck",
		"deck":        deck,
	})

	for _, conn := range s.clients {
		conn.WriteMessage(websocket.TextMessage, responseJSON)
	}
}

func (s *GameServer) handleUpdateCursor(data map[string]interface{}, userID string) {
	name, ok := data["name"].(string)
	if !ok {
		return
	}

	xPos, ok := data["x"].(float64)
	if !ok {
		return
	}

	yPos, ok := data["y"].(float64)
	if !ok {
		return
	}

	// Notify all users in the channel
	cursorJSON, _ := json.Marshal(map[string]interface{}{
		"type": "update_cursor",
		"from": userID,
		"name": name,
		"x":    xPos,
		"y":    yPos,
	})

	// send to all
	for str, conn := range s.clients {
		log.Println("sending to ", cursorJSON, str)
		conn.WriteMessage(websocket.TextMessage, cursorJSON)
	}

	responseJSON, _ := json.Marshal(map[string]interface{}{
		"type":    "success",
		"message": "sent cursor",
	})

	// send to sender
	s.clients[userID].WriteMessage(websocket.TextMessage, responseJSON)
}
