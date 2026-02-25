package game

// Current state of the Scoundrel game.
type Game struct {
	HP              int
	Dungeon         []*Card
	Room            []*Card
	LastDiscarded   *Card
	Weapon          *Card
	MonstersSlain   []*Card
	SkippedLastRoom bool
}

// Create a new game with the given deck.
func NewGame(d []*Card) *Game {
	return &Game{
		HP:              MaxHP,
		Dungeon:         d,
		Room:            make([]*Card, 0),
		LastDiscarded:   nil,
		Weapon:          nil,
		MonstersSlain:   make([]*Card, 0),
		SkippedLastRoom: false,
	}
}

// Create a new game with a shuffled Scoundrel deck.
func NewRandomGame() *Game {
	d := NewShuffledScoundrelDeck()
	return NewGame(d)
}

// Clear the existing room and deal a new one from the dungeon deck.
func (g *Game) DealRoom() {
	g.Room = g.Room[:0]
	for range min(CardsPerRoom, len(g.Dungeon)) {
		lastIdx := len(g.Dungeon) - 1
		g.Room = append(g.Room, g.Dungeon[lastIdx])
		g.Dungeon = g.Dungeon[:lastIdx]
	}
}

// Skip the current room, placing it at the bottom of the dungeon deck.
func (g *Game) SkipRoom() {
	var nonNils []*Card
	for _, c := range g.Room {
		if c != nil {
			nonNils = append(nonNils, c)
		}
	}
	g.Dungeon = append(nonNils, g.Dungeon...)
	g.SkippedLastRoom = true
	g.DealRoom()
}

// Check if enough actions have been taken to deal the next room.
func (g *Game) IsRoomDone() bool {
	nils := 0
	for _, c := range g.Room {
		if c == nil {
			nils += 1
			if nils >= ReqMovesPerRoom {
				return true
			}
		}
	}
	return false
}

// Add the monster card to the current slain ones.
func (g *Game) AddToSlain(monster *Card) {
	g.MonstersSlain = append(g.MonstersSlain, monster)
}

// Add player HP. Does not increase the HP after MaxHP.
func (g *Game) AddHP(p int) {
	if new := g.HP + p; new < MaxHP {
		g.HP = new
	} else {
		g.HP = MaxHP
	}
}

// Remove player HP. Does not decrease the HP after MinHP.
func (g *Game) RemoveHP(p int) {
	if new := g.HP - p; new > MinHP {
		g.HP = new
	} else {
		g.HP = MinHP
	}
}

// Use and discard a health potion card. roomIdx is the card's index in the current room.
func (g *Game) UseHealthPotion(roomIdx int) {
	potionCard := g.Room[roomIdx]
	g.AddHP(potionCard.IntRank())
	g.Room[roomIdx] = nil
	g.LastDiscarded = potionCard
}

// Take a new weapon card and discard the existing player hand, if any. roomIdx is the card's index in the current room.
func (g *Game) TakeWeapon(roomIdx int) {
	new := g.Room[roomIdx]

	if len(g.MonstersSlain) > 0 {
		g.MonstersSlain = g.MonstersSlain[:0]
	}
	if old := g.Weapon; old != nil {
		g.LastDiscarded = old
		g.Weapon = nil
	}

	g.Weapon = new
	g.Room[roomIdx] = nil
}

// Calculate the damage that would be taken by attacking the monster card. Does not affect HP or the room.
func (g *Game) CalculateDamage(monster *Card) int {
	fullDamage := monster.IntRank()

	if g.Weapon == nil {
		// no weapon -> full damage taken
		return fullDamage
	}

	// negative damage would add to HP!
	reducedDamage := max(0, fullDamage-g.Weapon.IntRank())

	if len(g.MonstersSlain) == 0 {
		// unused weapon has full strength
		return reducedDamage
	}

	lastSlain := g.MonstersSlain[len(g.MonstersSlain)-1]
	if lastSlain.RanksAbove(monster) {
		// last slain was of stronger than current one -> can use weapon
		return reducedDamage
	} else {
		// ... -> cannot use weapon
		return fullDamage
	}
}

// Take damage by the precalculated amount. roomIdx is the attacked monster's index in the current room.
func (g *Game) TakeDamage(damage, roomIdx int) {
	if damage < 0 {
		panic("damage cannot be negative")
	}
	monsterCard := g.Room[roomIdx]

	if g.Weapon != nil {
		g.AddToSlain(monsterCard)
	} else {
		g.LastDiscarded = monsterCard
	}
	g.Room[roomIdx] = nil
	g.RemoveHP(damage)
}
