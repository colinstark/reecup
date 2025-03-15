package game

type Player struct {
	UserID  string  `json:"userID"`
	Hand    []Stone `json:"hand"`
	HasMeld bool    `json:"hasMeld"`
}
