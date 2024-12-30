package game

type Player struct {
	Name  string
	Hand  []Stone
	IsOut bool
}

func NewPlayer(name string, deck Deck) Player {
	hand := []Stone{}
	for i := 0; i < 14; i++ {
		stone, err := deck.Draw()
		if err != nil {
			// handle error, but unlikely
		}
		hand = append(hand, stone)
	}

	return Player{
		Name:  name,
		Hand:  hand,
		IsOut: false,
	}
}
