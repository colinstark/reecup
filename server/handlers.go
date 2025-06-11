package server

import (
	"encoding/json"
	"log"
	"reecup/game"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func (s *GameServer) handleUpdateName(data map[string]any, userID userID) {
	name, ok := data["name"].(string)
	if !ok {
		s.handleError("login", "Invalid name", userID)
		return
	}

	s.userMutex.Lock()

	user := s.Users[userID]
	user.Name = name

	s.Users[userID] = user

	s.userMutex.Unlock()

	// Send success response
	s.Send(userID, "name_updated", map[string]any{
		"success": true,
		"userID":  userID,
		"name":    name,
	})
}

func (s *GameServer) handleLogin(data map[string]any, userID userID) {
	name, ok := data["name"].(string)
	if !ok {
		s.handleError("login", "Invalid name", userID)
		return
	}

	s.userMutex.Lock()
	s.Users[userID] = User{
		Name: name,
	}
	s.userMutex.Unlock()

	// Send success response
	s.Send(userID, "login", map[string]any{
		"success": true,
		"userID":  userID,
	})
}

func (s *GameServer) handleLogout(data map[string]any, userID userID) {
	s.userMutex.Lock()
	delete(s.Connections, userID)
	s.userMutex.Unlock()

	// Notify all other users that this user has logged out
	broadcastData := map[string]any{
		"instruction": "user_logout",
		"userID":      userID,
	}
	responseJSON, _ := json.Marshal(broadcastData)

	for id, conn := range s.Connections {
		if id != userID { // Don't send to the user who's logging out
			conn.WriteMessage(websocket.TextMessage, responseJSON)
		}
	}

	// Send confirmation to the logging out user
	s.handleError("logout", "success", userID)
}

func (s *GameServer) handleGetDeck(data map[string]any, userID string) {
	deck := game.CreateDeck()                   // Creates a new deck with stones
	hand := game.DrawForNewPlayer(userID, deck) // Draw initial hand for player

	responseJSON, _ := json.Marshal(map[string]any{
		"instruction": "new_deck",
		"hand":        hand,
		"remaining":   deck.GetCount(),
	})

	s.Connections[userID].WriteMessage(websocket.TextMessage, responseJSON)
}

func (s *GameServer) handleUpdateCursor(data map[string]any, userID string) {
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
	cursorJSON, _ := json.Marshal(map[string]any{
		"type": "update_cursor",
		"from": userID,
		"name": name,
		"x":    xPos,
		"y":    yPos,
	})

	// send to all
	for str, conn := range s.Connections {
		log.Println("sending to ", cursorJSON, str)
		conn.WriteMessage(websocket.TextMessage, cursorJSON)
	}

	responseJSON, _ := json.Marshal(map[string]any{
		"type":    "success",
		"message": "sent cursor",
	})

	// send to sender
	s.Connections[userID].WriteMessage(websocket.TextMessage, responseJSON)
}

func (s *GameServer) handleNewGame(userID string) {
	gameID := generateGameID()
	newGame := game.NewGame()

	s.gameMutex.Lock()
	s.Games[gameID] = newGame
	s.gameMutex.Unlock()

	// Notify all users about new game
	broadcastData := map[string]any{
		"instruction": "game_created",
		"game":        newGame,
	}
	responseJSON, _ := json.Marshal(broadcastData)
	for id, _ := range s.Users {
		s.Connections[id].WriteMessage(websocket.TextMessage, responseJSON)
	}
}

func (s *GameServer) handleListGames(data map[string]any, userID string) {
	s.gameMutex.RLock()
	gamesList := make([]game.Game, 0, len(s.Games))
	for _, g := range s.Games {
		if g.State == "waiting" { // only list games that are waiting for players
			gamesList = append(gamesList, g)
		}
	}
	s.gameMutex.RUnlock()

	responseJSON, _ := json.Marshal(map[string]any{
		"instruction": "games_list",
		"games":       gamesList,
	})
	s.Connections[userID].WriteMessage(websocket.TextMessage, responseJSON)
}

func (s *GameServer) handleJoinGame(data map[string]any, userID string) {
	gameID, ok := data["gameID"].(string)
	if !ok {
		s.handleError("join_game", "invalid game ID", userID)
		return
	}

	s.gameMutex.Lock()
	game, exists := s.Games[gameID]
	if !exists {
		s.gameMutex.Unlock()
		s.handleError("join_game", "game not found", userID)
		return
	}

	// Check if player is already in game
	for _, player := range game.Players {
		if player.Name == userID {
			s.gameMutex.Unlock()
			s.handleError("join_game", "already in game", userID)
			return
		}
	}

	// Add player to game
	// game.Players = append(game.Players, s.Users[userID]) @TODO fix this
	s.Games[gameID] = game
	s.gameMutex.Unlock()

	// Notify all users about player joining
	broadcastData := map[string]any{
		"instruction": "player_joined",
		"gameID":      gameID,
		"playerID":    userID,
		"game":        game,
	}

	s.Broadcast("banana", broadcastData)
}

func (s *GameServer) handleCancelGame(data map[string]any, userID string) {
	gameID, ok := data["gameID"].(string)
	if !ok {
		s.handleError("cancel_game", "invalid game ID", userID)
		return
	}

	s.gameMutex.Lock()
	game, exists := s.Games[gameID]
	if !exists {
		s.gameMutex.Unlock()
		s.handleError("cancel_game", "game not found", userID)
		return
	}

	// Check if game is in progress - if so, return error
	if game.State == "in_progress" {
		s.gameMutex.Unlock()
		s.handleError("cancel_game", "cannot cancel game in progress", userID)
		return
	}

	// Remove the game from the games map
	delete(s.Games, gameID)
	s.gameMutex.Unlock()

	// Notify all users about game cancellation
	broadcastData := map[string]any{
		"instruction": "game_cancelled",
		"gameID":      gameID,
	}

	s.Broadcast("game_cancelled", broadcastData)

	// Send success response to the user who cancelled
	s.Send(userID, "cancel_game", map[string]any{
		"success": true,
		"gameID":  gameID,
	})
}

func generateGameID() gameID {
	return "game_" + strconv.FormatInt(time.Now().UnixNano(), 6)
}

func generateUserID() userID {
	return "user_" + strconv.FormatInt(time.Now().UnixNano(), 6)
}
