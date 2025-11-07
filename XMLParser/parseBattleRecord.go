package main

import (
	"encoding/json"
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
	OtherFields [][]byte `xml:",any"`
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
	RoundNumber  int          `xml:"round"`
	PlayerData   PlayerData   `xml:"playerData"`
	ActionRecord ActionRecord `xml:"actionRecords"`
}

/* PlayerData and related data structures */
type PlayerData struct {
	ReactorCore         int                  `xml:"reactorCore"`
	Supply              int                  `xml:"supply"`
	PreRoundFightResult string               `xml:"preRoundFightResult"`
	NewUnitData         []NewUnitData        `xml:"units>NewUnitData"`
	CommanderSkills     []CommanderSkillData `xml:"commanderSkills>CommanderSkillData"`
	ActiveTechnologies  []ActiveTechnologies `xml:"activeTechnologies>UnitData"`
	Shop                Shop                 `xml:"shop"`
	Contraptions        Contraptions         `xml:"contraptions"`
	ContraptionIndex    []int                `xml:"contraptionIndex"`
	Blueprints          []int                `xml:"bluepints>int"` // There is a typo in the XML
	EnergyTowerSkils    []int                `xml:"energyTowerSkills>int"`
	TowerStrengthLevels []int                `xml:"towerStrengthLevels>int"`
	IsSpecialSupply     bool                 `xml:"IsSpecialSupply"`
}

type NewUnitData struct {
	Id          int      `xml:"id"`
	Index       int      `xml:"Index"`
	RoundCount  int      `xml:"RoundCount"`
	Durability  int      `xml:"Durability"`
	Exp         int      `xml:"Exp"`
	Level       int      `xml:"Level"`
	Position    Position `xml:"Position"`
	EquipmentId int      `xml:"EquipmentID"`
	IsRotate    bool     `xml:"IsRotate"`
	SellSupply  int      `xml:"SellSupply"`
}

type Position struct {
	X int `xml:"x"`
	Y int `xml:"y"`
}

type CommanderSkillData struct {
	Index        int  `xml:"index"`
	Id           int  `xml:"id"`
	IsActive     bool `xml:"isActive"`
	CoolingRound int  `xml:"coolingRound"`
}

type ActiveTechnologies struct {
	Id   int  `xml:"id,attr"`
	Tech Tech `xml:"techs>tech data"`
}

type Shop struct {
	UnlockedUnits []int `xml:"unlockedUnits>int"`
	LockedUnits   []int `xml:"lockedUnits>int"`
}

type Contraptions struct {
	ContraptionData []ContraptionData `xml:"ContraptionData"`
}

type ContraptionData struct {
	Index    int      `xml:"index"`
	Id       int      `xml:"id"`
	Position Position `xml:"position"`
}

/* ------------------------------ */

/* Action Records and related data structures */
type ActionRecord struct {
	MatchActionData []MatchActionData `xml:"MatchActionData"`
}

type MatchActionData struct {
	Type          string         `xml:"type,attr"`
	Time          int            `xml:"Time"`
	LocalTime     float32        `xml:"LocalTime"`
	ID            int            `xml:"ID,omitempty"`
	Index         int            `xml:"Index,omitempty"`
	UID           int            `xml:"UID,omitempty"`
	UIDX          int            `xml:"UIDX,omitempty"`
	TechID        int            `xml:"TechID,omitempty"`
	MoveUnitDatas []MoveUnitData `xml:"moveUnitDatas>MoveUnitData"`
}

type MoveUnitData struct {
	UnitId    int      `xml:"unitID"`
	UnitIndex int      `xml:"unitIndex"`
	Position  Position `xml:"position"`
	IsRotate  bool     `xml:"isRotate"`
}

/* ------------------------------ */

func BuildUnitTechLUT(record BattleRecord) map[int]map[int]string {
	lut := make(map[int]map[int]string)

	for _, player := range record.PlayerRecords {
		if player.Name != "Bot" {
			continue
		}
		for _, unit := range player.Data.UnitDatas {
			if unit.Techs != nil {
				for _, tech := range unit.Techs {
					if _, ok := lut[unit.Id]; !ok {
						lut[unit.Id] = make(map[int]string)
					}
					lut[unit.Id][tech.TechData] = "Tech"
				}
			}
		}
	}
	return lut
}

func saveParsedRecord(record BattleRecord, filename string) error {
	// Pretty-print JSON
	jsonData, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal record: %w", err)
	}

	return os.WriteFile(filename, jsonData, 0644)
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
	xmlFile, err := os.ReadFile("BattleRecords/vsAI_Insane_25-11-07-21-22-50.grbr")
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
		// break
	}

	for _, raw := range record.OtherFields {
		fmt.Println("Unknown tag content:\n", string(raw))
	}

	outputFile := "ParsedRecords/vsAI_Insane_25-11-07-21-22-50.json"
	if _, err := os.Stat(outputFile); err == nil {
		fmt.Println("File already exists, skipping: ", outputFile)
	} else {
		if err := saveParsedRecord(record, outputFile); err != nil {
			panic(err)
		}
		fmt.Println("Saved parsed XML as JSON to: ", outputFile)
	}

	lut := BuildUnitTechLUT(record)
	jsonData, _ := json.MarshalIndent(lut, "", "  ")
	os.WriteFile("unit_tech_lut.json", jsonData, 0644)
	fmt.Println("Exported LUT skeleton to unit_tech_lut.json")
}
