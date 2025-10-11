package main

import "fmt"

type TowerUpgrade struct {
	Name  string
	Level int
}

type UnitUpgrade struct {
	Name string
}

type Tower struct {
	Name      string
	Level     int
	Damage    float32
	Range     float32
	Position  [2]int // [x][y]
	Footprint [2]int // [x][y]
	Upgrades  []TowerUpgrade
}

type TargetType int

const (
	TargetGround TargetType = iota
	TargetAir
	TargetAll
)

type EnergyShield struct {
	Health float32
	Radius float32
}

type Unit struct {
	Name           string
	Level          int
	StartingLevel  int
	Cost           int
	UnlockCost     int
	Health         float32
	Speed          float32
	Attack         float32
	SplashRange    float32
	AttackInterval float32
	Range          float32
	TargetType     TargetType
	MaxTargets     int
	HealthRegen    float32
	Armor          float32
	Overshield     float32
	BubbleShield   *EnergyShield
	Position       [2]int // [x][y]
	Footprint      [2]int // [x][y]
	Quantity       [2]int // [cols][rows]
}

type StatBonus struct {
	EffectedUnitName string
	EffectedStatName string
	BonusValue       float64
}

type Specialist struct {
	Name           string
	StatBonuses    []StatBonus
	BonusUnit      Unit
	UnitSpawnRound int
}

type Augment struct {
	Name        string
	StatBonuses []StatBonus
}

type PlayerState struct {
	Health     int
	Money      int
	Income     int
	Towers     []Tower
	Units      []Unit
	Specialist Specialist
	Augments   []Augment
}

func main() {
	fmt.Println("Hello, World!")
}
