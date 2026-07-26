package spreadsheet

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
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

	// open/create log files
	validLog, err := os.OpenFile("spreadsheet/log_valid.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("error creating/opening log_valid.txt: %s\n", err.Error())
	}
	defer validLog.Close()

	matchLog, err := os.OpenFile("spreadsheet/log_match.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("error creating/opening log_match.txt: %s\n", err.Error())
	}
	defer matchLog.Close()

	typeLog, err := os.OpenFile("spreadsheet/log_type.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("error creating/opening log_type.txt: %s\n", err.Error())
	}
	defer typeLog.Close()

	// get valid rows
	validRows := make([]int, 0)
	scanner := bufio.NewScanner(validLog)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		rowID, err := strconv.Atoi(line)
		if err != nil {
			log.Println(err)
		}

		validRows = append(validRows, rowID)
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	fmt.Println(validRows)
	if len(validRows) == 0 {
		return
	}
	// update all organizers
	for i, o := range organizers {
		// skip valid row
		if len(validRows) > 0 && slices.Contains(validRows, i+2) {
			continue
		}

		err = updateOrganizerData(&o)
		if err != nil {
			errorMessage := fmt.Sprintf("error row %d: %s\n", i+2, err.Error())

			if strings.Contains(errorMessage, "match") {
				matchLog.WriteString(errorMessage)
			} else if strings.Contains(errorMessage, "type") {
				typeLog.WriteString(errorMessage)
			}

			log.Print(errorMessage)
			time.Sleep(time.Second*1 + time.Millisecond*100)

			continue
		}

		validLog.WriteString(strconv.Itoa(i+2) + "\n")
		fmt.Printf("Address: %s\nLat: %f\nLon: %f\n\n", o.Address, o.Latitude, o.Longitude)
		time.Sleep(time.Second*1 + time.Millisecond*100)
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
