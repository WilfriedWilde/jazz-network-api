package repository

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const organizerTableStatement = `
	CREATE TABLE IF NOT EXISTS organizers (
		id INTEGER PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		country TEXT NOT NULL,
		city TEXT NOT NULL,
		address TEXT NOT NULL,
		websites TEXT,
		emails TEXT,
		phone_numbers TEXT,
		contact_persons TEXT,
		schedule TEXT,
		notes TEXT,
		latitude REAL,
		longitude REAL
);`

func CreateOrganizer(db *sql.DB, o Organizer) error {
	_, err := db.Exec(`
		INSERT INTO organizers (type, name, country, city, address, websites, emails, phone_numbers, contact_persons, schedule, notes, latitude, longitude) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		o.Type, o.Name, o.Country, o.City, o.Address,
		strings.Join(o.Websites, ";"),
		strings.Join(o.Emails, ";"),
		strings.Join(o.PhoneNumbers, ";"),
		strings.Join(o.ContactPersons, ";"),
		o.Schedule, o.Notes, o.Latitude, o.Longitude,
	)

	if err != nil {
		return fmt.Errorf("error creating organizer: %s\n", err.Error())
	}

	fmt.Printf("Organizer created successfully:\n %+v\n\n", o)

	return nil
}

func ReadOrganizer(db *sql.DB, selectOpts Select) (*sql.Rows, error) {
	return db.Query(`SELECT ? ? FROM organizers ? ?`, formatDistinct(selectOpts.Distinct), strings.Join(selectOpts.ColumnList, ","), formatJoin("organizers", selectOpts.JoinTables), formatRowFilters(selectOpts.RowFilter))
}

func GetOrganizers() ([]Organizer, error) {
	file, err := os.Open("repository/jazz_network.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening jazz_network.csv: %s\n", err.Error())
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading jazz_network.csv: %s\n", err.Error())
	}

	organizers := make([]Organizer, len(records)-1)

	for rowID, row := range records[1:] {
		websites := getRowSlice(row[5])
		emails := getRowSlice(row[6])
		phoneNumbers := getRowSlice(row[7])
		contactPersons := getRowSlice(row[8])
		latitude, err := strconv.ParseFloat(row[11], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing latitude: %s\n", err.Error())
		}
		longitude, err := strconv.ParseFloat(row[11], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing longitude: %s\n", err.Error())
		}

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
			Latitude:       latitude,
			Longitude:      longitude,
		}
	}

	return organizers, nil
}
