package spreadsheet

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func UpdateSpreadSheet() {
	file, err := os.Open("spreadsheet/jazz_network.csv")
	if err != nil {
		log.Printf("error opening jazz_network.csv: %s\n", err.Error())
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Printf("error reading jazz_network.csv: %s\n", err.Error())
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
			Type:           strings.TrimSpace(strings.ToLower(row[0])),
			Name:           strings.TrimSpace(strings.ToLower(row[1])),
			Country:        strings.TrimSpace(strings.ToLower(row[2])),
			City:           strings.TrimSpace(strings.ToLower(row[3])),
			Address:        strings.TrimSpace(strings.ToLower(row[4])),
			Websites:       websites,
			Emails:         emails,
			PhoneNumbers:   phoneNumbers,
			ContactPersons: contactPersons,
			Schedule:       strings.TrimSpace(strings.ToLower(row[9])),
			Notes:          strings.TrimSpace(row[10]),
		}
	}

	errFile, err := os.OpenFile("log_error.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("error creating/opening log_error.txt: %s\n", err.Error())
	}

	for i, o := range organizers {
		err = updateOrganizerData(&o)
		if err != nil {
			errorMessage := fmt.Sprintf("error row %d: %s\n", i, err.Error())

			log.Print(errorMessage)
			errFile.WriteString(errorMessage)
			time.Sleep(time.Second*1 + time.Millisecond*200)

			continue
		}

		fmt.Printf("Address: %s\nLat: %f\nLon: %f\n\n", o.Address, o.Latitude, o.Longitude)
		time.Sleep(time.Second*1 + time.Millisecond*200)
	}
}

func getRowSlice(s string) []string {
	split := strings.Split(s, ";")
	rowSlice := make([]string, len(split))

	for i := range split {
		rowSlice[i] = strings.TrimSpace(strings.ToLower(split[i]))
	}

	return rowSlice
}
