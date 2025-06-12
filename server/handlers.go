package server

import (
	"encoding/json"
	"log"
	"math/rand"
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
	s.Broadcast("user_logout", map[string]any{
		"userID": userID,
	})

	// Send confirmation to the logging out user
	s.handleError("logout", "success", userID)
}

func (s *GameServer) handleGetDeck(data map[string]any, userID string) {
	deck := game.CreateDeck()                   // Creates a new deck with stones
	hand := game.DrawForNewPlayer(userID, deck) // Draw initial hand for player

	s.Send(userID, "new_deck", map[string]any{
		"hand":      hand,
		"remaining": deck.GetCount(),
	})
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
	s.Broadcast("game_created", map[string]any{
		"game": newGame,
	})
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

	s.Send(userID, "games_list", map[string]any{
		"games": gamesList,
	})
}

func (s *GameServer) handleJoinGame(data map[string]any, userID string) {
	gameID, ok := data["gameID"].(string)
	if !ok {
		s.handleError("join_game", "invalid game ID", userID)
		return
	}

	s.gameMutex.Lock()
	gameData, exists := s.Games[gameID]
	if !exists {
		s.gameMutex.Unlock()
		s.handleError("join_game", "game not found", userID)
		return
	}

	// Check if player is already in game
	for _, player := range gameData.Players {
		if player.Name == userID {
			s.gameMutex.Unlock()
			s.handleError("join_game", "already in game", userID)
			return
		}
	}

	// Add player to game
	userInfo := s.Users[userID]
	newPlayer := game.NewPlayer(userID, userInfo.Name)
	gameData.Players = append(gameData.Players, newPlayer)
	s.Games[gameID] = gameData
	s.gameMutex.Unlock()

	// Notify all users about player joining
	broadcastData := map[string]any{
		"instruction": "player_joined",
		"gameID":      gameID,
		"playerID":    userID,
		"game":        gameData,
	}

	s.BroadcastInGame(gameID, "player_joined", broadcastData)
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

	// Notify players in the game before deleting it
	broadcastData := map[string]any{
		"instruction": "game_cancelled",
		"gameID":      gameID,
	}

	s.BroadcastInGame(gameID, "game_cancelled", broadcastData)

	// Remove the game from the games map
	delete(s.Games, gameID)
	s.gameMutex.Unlock()

	// Send success response to the user who cancelled
	s.Send(userID, "cancel_game", map[string]any{
		"success": true,
		"gameID":  gameID,
	})
}

func (s *GameServer) handleStartGame(data map[string]any, userID string) {
	gameID, ok := data["gameID"].(string)
	if !ok {
		s.handleError("start_game", "invalid game ID", userID)
		return
	}

	s.gameMutex.Lock()
	game, exists := s.Games[gameID]
	if !exists {
		s.gameMutex.Unlock()
		s.handleError("start_game", "game not found", userID)
		return
	}

	// Check if game is in waiting state - if not, return error
	if game.State != "waiting" {
		s.gameMutex.Unlock()
		s.handleError("start_game", "game is not in waiting state", userID)
		return
	}

	// Check if there are players in the game
	if len(game.Players) == 0 {
		s.gameMutex.Unlock()
		s.handleError("start_game", "no players in game", userID)
		return
	}

	// Randomly select the first player
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(game.Players))
	game.CurrentPlayerTurn = &game.Players[randomIndex]

	// Update game state to in_progress
	game.State = "in_progress"
	game.StartedAt = time.Now()

	// Initialize the first turn
	game.StartNewTurn()

	s.Games[gameID] = game
	s.gameMutex.Unlock()

	// Notify all users about game starting
	broadcastData := map[string]any{
		"instruction":   "game_started",
		"gameID":        gameID,
		"game":          game,
		"currentPlayer": game.CurrentPlayerTurn,
	}

	s.BroadcastInGame(gameID, "game_started", broadcastData)

	// Send success response to the user who started the game
	s.Send(userID, "start_game", map[string]any{
		"success": true,
		"gameID":  gameID,
	})
}

func (s *GameServer) handleDrawStone(data map[string]any, userID string) {
	gameID, ok := data["gameID"].(string)
	if !ok {
		s.handleError("draw_stone", "invalid game ID", userID)
		return
	}

	s.gameMutex.Lock()
	game, exists := s.Games[gameID]
	if !exists {
		s.gameMutex.Unlock()
		s.handleError("draw_stone", "game not found", userID)
		return
	}

	// Check if game is in progress
	if game.State != "in_progress" {
		s.gameMutex.Unlock()
		s.handleError("draw_stone", "game is not in progress", userID)
		return
	}

	// Check if it's the player's turn
	if game.CurrentPlayerTurn == nil || game.CurrentPlayerTurn.ID != userID {
		s.gameMutex.Unlock()
		s.handleError("draw_stone", "not your turn", userID)
		return
	}

	// Draw a stone from the deck
	stone, err := game.Deck.Draw()
	if err != nil {
		s.gameMutex.Unlock()
		s.handleError("draw_stone", "no more stones in deck", userID)
		return
	}

	// Add stone to player's hand
	for i := range game.Players {
		if game.Players[i].ID == userID {
			game.Players[i].Hand = append(game.Players[i].Hand, stone)
			// Update the CurrentPlayerTurn pointer to point to the updated player
			game.CurrentPlayerTurn = &game.Players[i]
			break
		}
	}

	// Advance to next player's turn after drawing
	game.NextPlayerTurn()
	s.Games[gameID] = game
	s.gameMutex.Unlock()

	// Send the drawn stone to the player
	s.Send(userID, "stone_drawn", map[string]any{
		"success":        true,
		"stone":          stone,
		"remaining_deck": game.Deck.GetCount(),
		"currentPlayer":  game.CurrentPlayerTurn,
	})

	// Broadcast turn change to all players in the game
	s.BroadcastInGame(gameID, "turn_changed", map[string]any{
		"currentPlayer": game.CurrentPlayerTurn,
		"reason":        "stone_drawn",
	})

	log.Println("Player", userID, "drew stone:", stone, "- Next player:", game.CurrentPlayerTurn.ID)
}

func (s *GameServer) handleFinishTurn(data map[string]any, userID string) {
	gameID, ok := data["gameID"].(string)
	if !ok {
		s.handleError("finish_turn", "invalid game ID", userID)
		return
	}

	s.gameMutex.Lock()
	game, exists := s.Games[gameID]
	if !exists {
		s.gameMutex.Unlock()
		s.handleError("finish_turn", "game not found", userID)
		return
	}

	// Check if game is in progress
	if game.State != "in_progress" {
		s.gameMutex.Unlock()
		s.handleError("finish_turn", "game is not in progress", userID)
		return
	}

	// Check if it's the player's turn
	if game.CurrentPlayerTurn == nil || game.CurrentPlayerTurn.ID != userID {
		s.gameMutex.Unlock()
		s.handleError("finish_turn", "not your turn", userID)
		return
	}

	// First, validate that the pool is empty
	if !game.IsTempBoardPoolEmpty() {
		s.gameMutex.Unlock()
		log.Println("Player", userID, "attempted to finish turn with stones remaining in pool")
		s.handleError("finish_turn", "pool is not empty - all stones must be placed in sets", userID)
		return
	}

	// Validate the current turn's TempBoard
	if !game.CurrentTurn.TempBoard.AllSetsValid() {
		s.gameMutex.Unlock()
		log.Println("Player", userID, "attempted to finish turn with invalid board configuration")
		s.handleError("finish_turn", "invalid board configuration", userID)
		return
	}

	// Check if current player needs to achieve meld
	var playerNeedsMeld bool
	var playerIndex int
	for i := range game.Players {
		if game.Players[i].ID == userID {
			playerIndex = i
			playerNeedsMeld = !game.Players[i].HasMeld
			break
		}
	}

	// If player doesn't have meld, check if they achieved it this turn
	if playerNeedsMeld {
		if !game.CheckMeld() {
			s.gameMutex.Unlock()
			log.Println("Player", userID, "attempted to finish turn without achieving required meld (30+ points)")
			s.handleError("finish_turn", "must play stones worth more than 30 points for initial meld", userID)
			return
		}
		// Player achieved meld, update their status
		game.Players[playerIndex].HasMeld = true
		// Update the CurrentPlayerTurn pointer to reflect the change
		game.CurrentPlayerTurn = &game.Players[playerIndex]
		log.Println("Player", userID, "achieved initial meld!")
	}

	// If valid, replace the game's board with the TempBoard
	game.Board = game.CurrentTurn.TempBoard
	game.CurrentTurn.IsValid = true

	// Check if current player has won (empty hand)
	var winnerPlayerIndex = -1
	for i := range game.Players {
		if game.Players[i].ID == userID {
			if len(game.Players[i].Hand) == 0 {
				// Player has won!
				game.State = "finished"
				game.GameOver = true
				winnerPlayerIndex = i
				break
			}
		}
	}

	// If game is not over, advance to next player's turn
	if !game.GameOver {
		game.NextPlayerTurn()
	}

	s.Games[gameID] = game
	s.gameMutex.Unlock()

	if game.GameOver && winnerPlayerIndex >= 0 {
		// Broadcast game over message
		broadcastData := map[string]any{
			"instruction": "game_over",
			"gameID":      gameID,
			"winner":      game.Players[winnerPlayerIndex],
			"board":       game.Board,
		}

		s.BroadcastInGame(gameID, "game_over", broadcastData)

		// Send success response to the winner
		s.Send(userID, "finish_turn", map[string]any{
			"success":  true,
			"gameID":   gameID,
			"gameOver": true,
			"winner":   game.Players[winnerPlayerIndex],
			"board":    game.Board,
		})

		log.Println("Player", userID, "won the game!")
	} else {
		// Broadcast the turn change to all players in the game
		broadcastData := map[string]any{
			"instruction":   "turn_finished",
			"gameID":        gameID,
			"currentPlayer": game.CurrentPlayerTurn,
			"board":         game.Board,
		}

		s.BroadcastInGame(gameID, "turn_finished", broadcastData)

		// Send success response to the user who finished their turn
		s.Send(userID, "finish_turn", map[string]any{
			"success":       true,
			"gameID":        gameID,
			"currentPlayer": game.CurrentPlayerTurn,
			"board":         game.Board,
		})

		log.Println("Player", userID, "finished turn with valid board. Next player:", game.CurrentPlayerTurn.ID)
	}
}

func generateGameID() gameID {
	return "game_" + strconv.FormatInt(time.Now().UnixNano(), 6)
}

func generateUserID() userID {
	return "user_" + strconv.FormatInt(time.Now().UnixNano(), 6)
}
