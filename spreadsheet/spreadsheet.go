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
	// open, read and store records store raw csv
	rawFile, err := os.Open("spreadsheet/jazz_network_raw.csv")
	if err != nil {
		log.Printf("error opening jazz_network_raw.csv: %s\n", err.Error())
	}
	defer rawFile.Close()

	records, err := csv.NewReader(rawFile).ReadAll()
	if err != nil {
		log.Printf("error reading jazz_network_raw.csv: %s\n", err.Error())
	}

	organizers := make([]Organizer, len(records)-1)

	for rowID, row := range records[1:] {
		if row[3] == "" {
			fmt.Printf("no city for %s at line %d\n", row[1], rowID+2)
		}

		organizers[rowID] = Organizer{
			Type:           strings.TrimSpace(strings.ToLower(row[0])),
			Name:           strings.TrimSpace(strings.ToLower(row[1])),
			Country:        strings.TrimSpace(strings.ToLower(row[2])),
			City:           strings.TrimSpace(strings.ToLower(row[3])),
			Address:        strings.TrimSpace(strings.ToLower(row[4])),
			Websites:       strings.TrimSpace(strings.ToLower(row[5])),
			Emails:         strings.TrimSpace(strings.ToLower(row[6])),
			PhoneNumbers:   strings.TrimSpace(strings.ToLower(row[7])),
			ContactPersons: strings.TrimSpace(strings.ToLower(row[8])),
			Schedule:       strings.TrimSpace(strings.ToLower(row[9])),
			Notes:          strings.TrimSpace(row[10]),
		}
	}

	// open/create output csv
	outputFile, err := os.OpenFile("spreadsheet/jazz_network.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("error creating/opening jazz_network.csv: %s\n", err.Error())
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	headers := append(records[0], "LATITUDE", "LONGITUDE")
	if err := writer.Write(headers); err != nil {
		log.Printf("error writing headers to output file: %s", err)
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
		// // skip valid row
		// if len(validRows) > 0 && slices.Contains(validRows, i+2) {
		// 	continue
		// }

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

		if !slices.Contains(validRows, i+2) {
			validLog.WriteString(strconv.Itoa(i+2) + "\n")
		}

		fmt.Printf("Name: %s\nAddress: %s\nLat: %f\nLon: %f\n\n", o.Name, o.Address, o.Latitude, o.Longitude)

		// write new record to output file
		newRecord := []string{
			o.Type,
			o.Name,
			o.Country,
			o.City,
			o.Address,
			o.Websites,
			o.Emails,
			o.PhoneNumbers,
			o.ContactPersons,
			o.Schedule,
			o.Notes,
			strconv.FormatFloat(o.Latitude, 'g', -1, 64),
			strconv.FormatFloat(o.Longitude, 'g', -1, 64),
		}

		if err := writer.Write(newRecord); err != nil {
			log.Printf("error writing new record to output file: %s", err)
		}

		time.Sleep(time.Second*1 + time.Millisecond*100)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Printf("error writing to jazz_network.csv: %s", err.Error())
	}
}
