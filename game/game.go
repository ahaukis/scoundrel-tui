package game

// Current state of the Scoundrel game.
type Game struct {
	HP                     int
	Dungeon                []*Card
	Room                   []*Card
	LastDiscarded          *Card
	Weapon                 *Card
	MonstersSlain          []*Card
	skippedLastRoom        bool
	usedHealthPotionInRoom bool
}

// Create a new game with the given deck.
func NewGame(d []*Card) *Game {
	return &Game{
		HP:            MaxHP,
		Dungeon:       d,
		Room:          make([]*Card, 0),
		LastDiscarded: nil,
		Weapon:        nil,
		MonstersSlain: make([]*Card, 0),
	}
}

// Create a new game with a shuffled Scoundrel deck.
func NewRandomGame() *Game {
	d := NewShuffledScoundrelDeck()
	return NewGame(d)
}

func (g *Game) Lost() bool {
	return g.HP <= 0
}

func (g *Game) Won() bool {
	return g.HP > 0 && len(g.nonNilRoomCards()) == 0 && len(g.Dungeon) == 0
}

// Clear the existing room and deal a new one from the dungeon deck.
func (g *Game) DealRoom() {
	for _, c := range g.Room {
		if c != nil {
			g.LastDiscarded = c
		}
	}
	g.Room = g.Room[:0]

	for range min(CardsPerRoom, len(g.Dungeon)) {
		lastIdx := len(g.Dungeon) - 1
		g.Room = append(g.Room, g.Dungeon[lastIdx])
		g.Dungeon = g.Dungeon[:lastIdx]
	}

	if g.skippedLastRoom {
		g.skippedLastRoom = false
	}
	if g.usedHealthPotionInRoom {
		g.usedHealthPotionInRoom = false
	}
}

func (g *Game) nonNilRoomCards() []*Card {
	var nonNils []*Card
	for _, c := range g.Room {
		if c != nil {
			nonNils = append(nonNils, c)
		}
	}
	return nonNils
}

// Skip the current room, placing it at the bottom of the dungeon deck.
func (g *Game) SkipRoom() {
	// cannot skip 2 rooms in a row
	if g.skippedLastRoom {
		return
	}
	// cannot skip after already enganing a room...
	if nonNils := g.nonNilRoomCards(); len(nonNils) < CardsPerRoom && len(g.Dungeon) > 0 {
		// ...unless the only cards left are health potions and cannot be consumed
		allHealthPotions := true
		for _, c := range nonNils {
			if !c.IsHealthPotion() {
				allHealthPotions = false
				break
			}
		}
		if !(allHealthPotions && g.usedHealthPotionInRoom) {
			return
		}
	}

	g.Dungeon = append(g.nonNilRoomCards(), g.Dungeon...)
	g.DealRoom()
	g.skippedLastRoom = true
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

func (g *Game) MakeRoomAction(roomIdx int, useWeapon bool) {
	c := g.Room[roomIdx]
	if c == nil {
		return
	}
	switch {
	case c.IsWeapon():
		g.TakeWeapon(roomIdx)
	case c.IsHealthPotion():
		g.UseHealthPotion(roomIdx)
	case c.IsMonster():
		g.AttackMonster(roomIdx, useWeapon)
	}
}

// Use and discard a health potion card. roomIdx is the card's index in the current room.
func (g *Game) UseHealthPotion(roomIdx int) {
	if g.usedHealthPotionInRoom {
		// can only ue 1 health potion per room
		return
	}
	potionCard := g.Room[roomIdx]
	g.AddHP(potionCard.IntRank())
	g.Room[roomIdx] = nil
	g.LastDiscarded = potionCard
	g.usedHealthPotionInRoom = true
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

func (g *Game) AttackMonster(roomIdx int, useWeapon bool) {
	m := g.Room[roomIdx]
	dmg, weaponFailed := g.calculateDamage(m, useWeapon)

	if dmg < 0 {
		panic("damage cannot be negative")
	}

	if g.Weapon != nil && useWeapon && !weaponFailed {
		g.AddToSlain(m)
	} else {
		g.LastDiscarded = m
	}
	g.Room[roomIdx] = nil
	g.RemoveHP(dmg)
}

// Calculate the damage that would be taken by attacking the monster card. Does not affect HP or the room.
func (g *Game) calculateDamage(monster *Card, useWeapon bool) (int, bool) {
	fullDamage := monster.IntRank()

	if g.Weapon == nil || !useWeapon {
		// no weapon -> full damage taken
		return fullDamage, false
	}

	// negative damage would add to HP!
	reducedDamage := max(0, fullDamage-g.Weapon.IntRank())

	if len(g.MonstersSlain) == 0 {
		// unused weapon has full strength
		return reducedDamage, false
	}

	lastSlain := g.MonstersSlain[len(g.MonstersSlain)-1]
	if lastSlain.RanksAbove(monster) {
		// last slain was of stronger than current one -> can use weapon
		return reducedDamage, false
	} else {
		// ... -> cannot use weapon
		return fullDamage, true
	}
}
