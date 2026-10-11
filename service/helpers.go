package service

import (
	"jazz-network-api/model"
	"slices"
	"time"
)

func getNumTourStops(tour model.Tour) int {
	start, _ := time.Parse(time.DateOnly, tour.Start)
	end, _ := time.Parse(time.DateOnly, tour.End)

	return int(end.Sub(start).Hours()/24) + 1
}

func getMissingDates(existing []string, start, end string) ([]string, error) {
	startDate, err := time.Parse(time.DateOnly, start)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(time.DateOnly, end)
	if err != nil {
		return nil, err
	}

	lookupDate := make(map[string]struct{}, len(existing))
	for _, d := range existing {
		lookupDate[d] = struct{}{}
	}

	var missing []string
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		date := d.Format(time.DateOnly)

		if _, found := lookupDate[date]; !found {
			missing = append(missing, date)
		}
	}

	slices.Sort(missing)

	return missing, nil
}
