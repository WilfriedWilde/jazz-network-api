package countries

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/biter777/countries"
)

const fileLocation = "countries/countries.json"

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func GetCountriesJSON() (*os.File, error) {
	var file *os.File

	if _, err := os.Stat(fileLocation); errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(fileLocation)
		if err != nil {
			return nil, fmt.Errorf("create json file: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("stat json file: %w", err)
	} else {
		file, err = os.Open(fileLocation)
		if err != nil {
			return nil, fmt.Errorf("open json file: %w", err)
		}

		return file, nil
	}

	allCountries := countries.All()

	result := make([]Country, len(allCountries))
	for i, c := range allCountries {
		result[i] = Country{
			Name: c.Info().Name,
			Code: c.Alpha2(),
		}
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return nil, fmt.Errorf("encode countries: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek json file: %w", err)
	}

	return file, nil
}

func GetCountries() (map[string]Country, error) {
	file, err := GetCountriesJSON()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var countries []Country

	if err := json.NewDecoder(file).Decode(&countries); err != nil {
		return nil, fmt.Errorf("decode countries: %w", err)
	}

	byCode := make(map[string]Country, len(countries))
	for _, c := range countries {
		byCode[c.Code] = c
	}

	return byCode, nil
}
