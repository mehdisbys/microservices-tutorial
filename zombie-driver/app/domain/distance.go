package domain

import (
	"math"
)

const earthRadius = float64(6371)

type DistanceEstimator interface {
	Distance(lat, long Coordinates) float64
	SumDistanceList(list []Coordinates) float64
}

type HaversineDistance struct{}

func (h HaversineDistance) Distance(from, to Coordinates) (distance float64) {
	toLat, toLong := to.GetLatitude(), to.GetLongitude()
	fromLat, fromLong := from.GetLatitude(), from.GetLongitude()

	var deltaLat = (toLat - fromLat) * (math.Pi / 180)

	var deltaLon = (toLong - fromLong) * (math.Pi / 180)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(fromLat*(math.Pi/180))*math.Cos(toLat*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance = earthRadius * c

	return math.Round(distance*1000) / 1000
}

func (h HaversineDistance) SumDistanceList(list []Coordinates) float64 {
	res := 0.0
	for i := range list {
		if len(list) > i+1 {
			res += h.Distance(list[i], list[i+1])
		}
	}

	return math.Round(res*1000) / 1000
}
