package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type BattleRecord struct {
	XMLName       xml.Name       `xml:"BattleRecord"`
	PlayerRecords []PlayerRecord `xml:"playerRecords>PlayerRecord"`
	// Add other fields as needed
}

type PlayerRecord struct {
	Name        string              `xml:"name"`
	Seed        int                 `xml:"seed"`
	Ad          int                 `xml:"ad"`
	Data        Data                `xml:"data"`
	RoundRecord []PlayerRoundRecord `xml:"playerRoundRecords>PlayerRoundRecord"`
	// Add other fields as needed
}

type Data struct {
	ReactorCore              int        `xml:"reactorCore"`
	MaxReactorCore           int        `xml:"maxReactorCore"`
	MaxRoundSupply           int        `xml:"maxRoundSupply"`
	FirstRoundSupply         int        `xml:"firstRoundSupply"`
	RoundSupplyIncreaseValue int        `xml:"roundSupplyIncreaseValue"`
	Team                     int        `xml:"team"`
	IsLeader                 bool       `xml:"isLeader"`
	Type                     string     `xml:"type"`
	UnitDatas                []UnitData `xml:"unitDatas>unitData"`
}

type UnitData struct {
	Id           int    `xml:"id"`
	Techs        []Tech `xml:"techs>tech"`
	UnlockedTech []Tech `xml:"unlockedTech>tech"`
}

type Tech struct {
	TechData int `xml:"data,attr"` // Extracts an int that likely correlates to a tech id in game, TODO: Create LUT mapping techs to id's
}

type PlayerRoundRecord struct {
	RoundNumber int        `xml:"round"`
	PlayerData  PlayerData `xml:"playerData"`
}

type PlayerData struct {
	ReactorCore         int           `xml:"reactorCore"`
	Supply              int           `xml:"supply"`
	PreRoundFightResult string        `xml:"preRoundFightResult"`
	NewUnitData         []NewUnitData `xml:"units>NewUnitData"`
}

type NewUnitData struct {
	Id          int
	Index       int
	RoundCount  int
	Durability  int
	Exp         int
	Level       int
	Position    Position
	EquipmentId int
	IsRotate    bool
	SellSupply  int
}

type Position struct {
	X int
	Y int
}

func printFields(v reflect.Value, indent string) {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			fmt.Printf("%s%s: \n", indent, field.Name)
			printFields(value, indent+"  ")
		}
	case reflect.Slice:
		fmt.Printf("[\n")
		for i := 0; i < v.Len(); i++ {
			printFields(v.Index(i), indent+"  ")
		}
		fmt.Printf("%s]\n", indent)
	default:
		fmt.Printf("%s%v\n", indent, v.Interface())
	}
}

func main() {
	xmlFile, err := os.ReadFile("battle_record.grbr")
	if err != nil {
		panic(err)
	}

	fileContent := string(xmlFile)
	idx := strings.Index(fileContent, "<?xml")
	if idx == -1 {
		panic("No XML content found")
	}
	xmlData := fileContent[idx:]

	var record BattleRecord
	err = xml.Unmarshal([]byte(xmlData), &record)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed %d player records\n", len(record.PlayerRecords))
	// Prints results of all players iteratively
	for _, player := range record.PlayerRecords {
		if player.Name != "" {
			fmt.Println("Player Record:")
			printFields(reflect.ValueOf((player)), "  ")
			fmt.Println("--------------------------------------------------")
		}
	}
}
