package spreadsheet

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("jazz_venues.csv")
	if err != nil {
		log.Fatal("error opening jazz_venues.csv: %s", err.Error())
	}

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal("error reading jazz_venues.csv: %s", err.Error())
	}

	fmt.Println(records[0])
}
