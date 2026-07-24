package spreadsheet

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

type Organizer struct {
	Type           string
	Name           string
	Country        string
	City           string
	Address        string
	Websites       []string
	Emails         []string
	PhoneNumbers   []string
	ContactPersons []string
	Schedule       string
	Notes          string
	Latitude       float64
	Longitude      float64
	Status         bool
}

func UpdateSpreadSheet() {
	file, err := os.Open("spreadsheet/jazz_network.csv")
	if err != nil {
		log.Printf("error opening jazz_network.csv: %s", err.Error())
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Printf("error reading jazz_network.csv: %s", err.Error())
	}

	organizers := make([]Organizer, len(records)-1)

	for rowID, row := range records[1:] {
		if row[3] == "" {
			fmt.Printf("no city for %s at line %d\n", row[1], rowID+2)
		}

		websites := getRowSlice(row[5])
		emails := getRowSlice(row[6])
		phoneNumbers := getRowSlice(row[7])
		contactPersons := getRowSlice(row[8])

		organizers[rowID] = Organizer{
			Type:           strings.TrimSpace(row[0]),
			Name:           strings.TrimSpace(row[1]),
			Country:        strings.TrimSpace(row[2]),
			City:           strings.TrimSpace(row[3]),
			Address:        strings.TrimSpace(row[4]),
			Websites:       websites,
			Emails:         emails,
			PhoneNumbers:   phoneNumbers,
			ContactPersons: contactPersons,
			Schedule:       strings.TrimSpace(row[9]),
			Notes:          strings.TrimSpace(row[10]),
		}
	}

	fmt.Printf("%+v\n", organizers[2])
}

func getRowSlice(s string) []string {
	split := strings.Split(s, ";")
	rowSlice := make([]string, len(split))

	for i := range split {
		rowSlice[i] = strings.TrimSpace(split[i])
	}

	return rowSlice
}
