package game

type Board struct {
	Sets []Set   `json:"sets"`
	Pool []Stone `json:"pool"`
}


func (b Board) AllSetsValid() bool {
	if len(b.Pool) != 0 {
		return false
	}
	for _, Set := range b.Sets {
		if Set.Validate() == false {
			return false
		}
	}
	return true
}


// func (b *Board) AddChain(sequential bool, stones []Stone) error {
// 	chain := Set{
// 		sequential: sequential,
// 		Stones:     stones,
// 	}
// 	if chain.SetValid() {
// 		b.Sets = append(b.Sets, chain)
// 		return nil
// 	} else {
// 		return errors.New("not a valid chain")
// 	}
// }

// func (b *Board) PopSet(idx int) Set {
// 	// pop chain for modification at idx
// 	//return chain

// 	return Set{}
// }
