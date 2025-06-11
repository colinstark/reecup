package game

type Player struct {
	ID      string  `json:"userID"`
	Name    string  `json:"userName"`
	Hand    []Stone `json:"hand"`
	HasMeld bool    `json:"hasMeld"`
}
