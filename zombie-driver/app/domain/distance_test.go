package domain

import (
	"testing"

	"github.com/go-test/deep"
)

func TestHaversineDistance(t *testing.T) {
	haversine := HaversineDistance{}

	tests := []struct {
		name     string
		from     Coordinates
		to       Coordinates
		d        DistanceEstimator
		expected float64
	}{
		{
			name:     "Distance Paris - London",
			from:     Coordinates{Lat: 48.8566, Long: 2.3522},
			to:       Coordinates{Lat: 51.5074, Long: 0.1278},
			d:        haversine,
			expected: 334.576,
		},
		{
			name:     "Distance is zero",
			from:     Coordinates{Lat: 48.8566, Long: 2.3522},
			to:       Coordinates{Lat: 48.8566, Long: 2.3522},
			d:        haversine,
			expected: 0.0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := test.d.Distance(test.from, test.to)

			if diff := deep.Equal(res, test.expected); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestHaversineDistanceSum(t *testing.T) {
	haversine := HaversineDistance{}
	tests := []struct {
		name     string
		list     []Coordinates
		d        DistanceEstimator
		expected float64
	}{
		{
			name: "Odd length",
			list: []Coordinates{
				{Lat: 48.8566, Long: 2.3522},   // Paris
				{Lat: 51.5074, Long: 0.1278},   // London
				{Lat: 35.6762, Long: 139.6503}, // Tokyo
			},
			d:        haversine,
			expected: 9883.81,
		},
		{
			name: "Even length",
			list: []Coordinates{
				{Lat: 48.8566, Long: 2.3522},   // Paris
				{Lat: 51.5074, Long: 0.1278},   // London
				{Lat: 35.6762, Long: 139.6503}, // Tokyo
				{Lat: 41.3851, Long: 2.1734},   // Barcelona
			},
			d:        haversine,
			expected: 20296.948,
		},
		{
			name:     "Empty list",
			list:     []Coordinates{},
			d:        haversine,
			expected: 0.0,
		},
		{
			name: "Same coordinates - immobility",
			list: []Coordinates{
				{Lat: 48.8566, Long: 2.3522},
				{Lat: 48.8566, Long: 2.3522},
				{Lat: 48.8566, Long: 2.3522},
			},
			d:        haversine,
			expected: 0.0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := test.d.SumDistanceList(test.list)

			if diff := deep.Equal(res, test.expected); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestIsZombie(t *testing.T) {
	tests := []struct {
		name            string
		driverDistance  float64
		minimumDistance float64
		expected        bool
	}{
		{
			name:            "Is zombie",
			driverDistance:  100.0,
			minimumDistance: 500.0,
			expected:        true,
		},
		{
			name:            "Is not zombie",
			driverDistance:  1000.0,
			minimumDistance: 500.0,
			expected:        false,
		},
		{
			name:            "Is not zombie - zero minimum",
			driverDistance:  1000.0,
			minimumDistance: 0.0,
			expected:        false,
		},
		{
			name:            "Is not zombie - zero minimum - zero distance",
			driverDistance:  0.0,
			minimumDistance: 0.0,
			expected:        false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			z := ZombieDetector{MinimumDistance: test.minimumDistance}

			res := z.IsZombie(test.driverDistance)

			if diff := deep.Equal(res, test.expected); diff != nil {
				t.Error(diff)
			}
		})
	}
}
