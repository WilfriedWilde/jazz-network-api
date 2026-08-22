package database

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"jazz-network-api/model"
	r "jazz-network-api/repository"
)

func seed(db *sql.DB) error {
	organizers, err := getOrganizersFromCSV()
	if err != nil {
		return fmt.Errorf("error getting organizers from csv file: %w", err)
	}

	for _, o := range organizers {
		if _, err := r.NewOrganizerRepository(db).Create(context.Background(), o); err != nil {
			return err
		}
	}

	return nil
}

func getOrganizersFromCSV() ([]model.Organizer, error) {
	file, err := os.Open("spreadsheet/jazz_network.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening jazz_network.csv: %w", err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading jazz_network.csv: %w", err)
	}

	organizers := make([]model.Organizer, len(records)-1)

	for rowID, row := range records[1:] {
		latitude, err := strconv.ParseFloat(row[11], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing latitude: %w", err)
		}
		longitude, err := strconv.ParseFloat(row[12], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing longitude: %w", err)
		}

		organizers[rowID] = model.Organizer{
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
			Latitude:       latitude,
			Longitude:      longitude,
		}
	}

	return organizers, nil
}
