package game

type Player struct {
	ID      string  `json:"userID"`
	Name    string  `json:"userName"`
	Hand    []Stone `json:"hand"`
	HasMeld bool    `json:"hasMeld"`
}

func NewPlayer(id, name string) Player {
	return Player{
		ID:      id,
		Name:    name,
		Hand:    []Stone{},
		HasMeld: false,
	}
}
