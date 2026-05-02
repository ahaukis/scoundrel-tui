package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	g := NewRandomGame()

	assert.Equalf(t, MaxHP, g.HP, "HP = %d; want %d", g.HP, MaxHP)
	assert.Falsef(t, g.Lost() || g.Won(), "Game begun in loss or win state: lost=%t, won=%t", g.Lost(), g.Won())

}

func TestDealRoom(t *testing.T) {
	g := NewRandomGame()
	assert.Empty(t, g.Room)

	initialDungeonSize := len(g.Dungeon)
	g.DealRoom()

	assert.Lenf(t, g.Room, CardsPerRoom, "%d room cards; want %d", len(g.Room), CardsPerRoom)
	assert.Lenf(t, g.Dungeon, initialDungeonSize-CardsPerRoom, "%d dungeon cards; want %d, %d", len(g.Dungeon), initialDungeonSize-CardsPerRoom, initialDungeonSize)
}

func TestSkipRoom(t *testing.T) {
	g := NewRandomGame()
	g.DealRoom()
	assert.True(t, g.CanSkipRoom())
	g.SkipRoom()
	assert.False(t, g.CanSkipRoom())
}

func TestUseHealthPotion(t *testing.T) {
	g := NewRandomGame()
	g.HP = MaxHP - 2
	potion := NewCard(Rank(2), Hearts)
	g.Room = append(g.Room, potion)
	g.UseHealthPotion(0)

	assert.Equal(t, MaxHP, g.HP)
	assert.Equal(t, potion, g.LastDiscarded)
	assert.Nil(t, g.Room[0])
	assert.True(t, g.usedHealthPotionInRoom)
}

func TestUseHealthPotionFails(t *testing.T) {
	g := NewRandomGame()
	g.HP = MaxHP - 2
	potion := NewCard(Rank(2), Hearts)
	g.Room = append(g.Room, potion)
	g.usedHealthPotionInRoom = true
	g.UseHealthPotion(0)

	assert.Equal(t, MaxHP-2, g.HP)
	assert.Equal(t, potion, g.Room[0])
	assert.True(t, g.usedHealthPotionInRoom)
}

func TestTakeWeapon(t *testing.T) {
	g := NewRandomGame()
	oldWeapon := NewCard(Rank(5), Diamonds)
	slain := NewCard(Rank(3), Spades)
	g.Weapon = oldWeapon
	g.MonstersSlain = append(g.MonstersSlain, slain)

	newWeapon := NewCard(Rank(10), Diamonds)
	g.Room = append(g.Room, newWeapon)
	g.TakeWeapon(0)

	assert.Nil(t, g.Room[0])
	assert.Equal(t, newWeapon, g.Weapon)
	assert.Equal(t, oldWeapon, g.LastDiscarded)
	assert.Empty(t, g.MonstersSlain)
}

func TestAttackMonsterWithWeapon(t *testing.T) {
	g := NewRandomGame()
	m := NewCard(Rank(5), Clubs)
	g.Weapon = NewCard(Rank(3), Diamonds)
	g.Room = append(g.Room, m)
	g.AttackMonster(0, true)

	assert.Equal(t, m, g.MonstersSlain[0])
	assert.Equal(t, MaxHP-2, g.HP)
}

func TestAttackMonsterNoWeapon(t *testing.T) {
	g := NewRandomGame()
	m := NewCard(Rank(5), Clubs)
	g.Room = append(g.Room, m)
	g.AttackMonster(0, true)

	assert.Equal(t, m, g.LastDiscarded)
	assert.Equal(t, MaxHP-5, g.HP)
}

func TestCalculateDamageFreshWeapon(t *testing.T) {
	var tests = []struct {
		monsterRank int
		weaponRank  int
		useWeapon   bool
		wantDamage  int
	}{
		{5, 5, true, 0},
		{5, 5, false, 5},
		{10, 5, true, 5},
		{10, 5, false, 10},
		{5, 10, true, 0},
		{5, 10, false, 5},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%d-%t", tt.monsterRank, tt.weaponRank, tt.useWeapon), func(t *testing.T) {
			g := NewRandomGame()
			m := NewCard(Rank(tt.monsterRank), Clubs)
			g.Weapon = NewCard(Rank(tt.weaponRank), Diamonds)

			damage, weaponFailed := g.calculateDamage(m, tt.useWeapon)

			assert.Equal(t, tt.wantDamage, damage)
			assert.False(t, weaponFailed)
		})
	}
}

func TestCalculateDamageUsedWeapon(t *testing.T) {
	var tests = []struct {
		monsterRank      int
		lastSlainRank    int
		weaponRank       int
		useWeapon        bool
		wantDamage       int
		wantWeaponFailed bool
	}{
		{5, 7, 10, true, 0, false},
		{5, 7, 10, false, 5, false},
		{5, 3, 10, true, 5, true},
		{5, 3, 10, false, 5, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%d-%d-%t", tt.monsterRank, tt.lastSlainRank, tt.weaponRank, tt.useWeapon), func(t *testing.T) {
			g := NewRandomGame()
			m := NewCard(Rank(tt.monsterRank), Clubs)
			g.Weapon = NewCard(Rank(tt.weaponRank), Diamonds)
			g.MonstersSlain = append(g.MonstersSlain, NewCard(Rank(tt.lastSlainRank), Clubs))

			damage, weaponFailed := g.calculateDamage(m, tt.useWeapon)

			assert.Equal(t, tt.wantDamage, damage)
			assert.Equal(t, tt.wantWeaponFailed, weaponFailed)
		})
	}
}
