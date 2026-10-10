package service

import (
	"jazz-network-api/model"
	"time"
)

func getNumTourStops(tour model.Tour) int {
	start, _ := time.Parse(time.DateOnly, tour.Start)
	end, _ := time.Parse(time.DateOnly, tour.End)

	return int(end.Sub(start).Hours()/24) + 1
}
