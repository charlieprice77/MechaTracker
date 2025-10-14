package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type BattleRecord struct {
	XMLName       xml.Name       `xml:"BattleRecord"`
	PlayerRecords []PlayerRecord `xml:"playerRecords>PlayerRecord"`
	// Add other fields as needed
}

type PlayerRecord struct {
	Name string `xml:"name"`
	// Add other fields as needed
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
	for _, player := range record.PlayerRecords {
		fmt.Println("Player Name:", player.Name)
	}
}
