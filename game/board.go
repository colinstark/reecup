package game

type Chain struct {
	sequential bool
	Stones     []Stone
}

type Board struct {
	Chains        []Chain
	floatingJoker bool
}

func (b Board) AllChainsValid() bool {
	for _, chain := range b.Chains {
		if chain.ChainValid() == false {
			return false
		}
	}
	return true
}

func (c Chain) ChainValid() bool {
	if c.sequential == true {
		// make sure order is correct, including 13 -> 1
		// make sure only one color
	} else {
		// make sure all faces are equal or joker
		// make sure only different colors
	}
	return false
}

func (b *Board) AddChain(sequential bool, stones []Stone) error {
	// check if valid chain
	b.Chains = append(b.Chains, Chain{
		sequential: sequential,
		Stones:     stones,
	})
	return nil
}
