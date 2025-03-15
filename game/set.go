package game

type Set struct {
	Stones []Stone `json:"stones"`
	Valid  bool    `json:"valid"`
}

func (s *Set) AddStone(idx int, stone Stone) {
	if idx < 0 || idx > len(s.Stones) {
		return
	}

	newStones := make([]Stone, 0, len(s.Stones)+1)
	newStones = append(newStones, s.Stones[:idx]...)
	newStones = append(newStones, stone)
	newStones = append(newStones, s.Stones[idx+1:]...)

	s.Stones = newStones
}

func (s *Set) RemoveStone(idx int) Stone {
	if idx < 0 || idx > len(s.Stones) {
		return Stone{}
	}
	newStones := make([]Stone, 0, len(s.Stones)-1)
	newStones = append(newStones, s.Stones[:idx]...)
	newStones = append(newStones, s.Stones[idx+1:]...)

	s.Stones = newStones
	return s.Stones[idx]
}

func (s *Set) Split(idx int) Set {
	left := s.Stones[:idx]
	right := s.Stones[idx+1:]
	s.Stones = left

	return Set{Stones: right}
}

func (s *Set) Move(old int, new int) {
	stone := s.RemoveStone(old)
	s.AddStone(new, stone)
}

func (s *Set) Validate() bool {
	// can't be smaller than 3
	if len(s.Stones) < 3 {
		return false
	}

	// check if has only one color
	isSequential := true
	setColor := s.Stones[0].Color
	if s.Stones[0].Joker {
		setColor = s.Stones[1].Color
	}
	for _, stone := range s.Stones {
		if stone.Color != setColor && stone.Joker == false {
			isSequential = false
		}
	}

	if isSequential == true {
		// make sure order is correct, including 13 -> 1
		for i := 1; i < len(s.Stones); i++ {
			if s.Stones[i].Joker == false && s.Stones[i-1].Joker == false {
				if s.Stones[i].Face != s.Stones[i-1].Face+1 {
					if s.Stones[i].Face == 1 {
						// Special case: if the previous value is 1, the next should be 13
						if s.Stones[i-1].Face != 13 {
							return false
						}
					} else {
						return false
					}
				}
			}

		}

		return true
	} else {
		// make sure all faces are equal or joker
		setFace := s.Stones[0].Face
		for _, stone := range s.Stones {
			if stone.Face != setFace && stone.Joker == false {
				return false
			}
		}

		// make sure only different colors
		colorSet := make(map[string]bool)

		for _, stone := range s.Stones {
			if stone.Joker == false {
				if _, ok := colorSet[stone.Color]; ok {
					return false // color already found
				}
				colorSet[stone.Color] = true
			}
		}

		return true
	}
}
