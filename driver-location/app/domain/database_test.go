// +build integration

package domain

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-test/deep"
	"github.com/heetch/MehdiSouilhed-technical-test/common"
)

// Make sure a redis instance is running or these tests will fail

func init() {
	database := NewInMemoryDB(redis.NewClient(&redis.Options{}))

	_, err := database.client.Ping().Result()

	if err != nil {
		panic(err)
	}
	database.client.Del("1") // clean data for driverID 1
}

func TestDatabase(t *testing.T) {
	now := time.Now().UTC()
	n := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	tests := []struct {
		name        string
		driverID    string
		coordinates []Coordinates
		expected    *[]Coordinates
	}{
		{
			name:     "Set/Fetch one coordinates",
			driverID: "1",
			coordinates: []Coordinates{
				{
					Lat:       1,
					Long:      2,
					UpdatedAt: common.Timestamp{Time: n.UTC()},
				},
			},
			expected: &[]Coordinates{
				{
					Lat:       1,
					Long:      2,
					UpdatedAt: common.Timestamp{Time: n.UTC()},
				},
			},
		},
		{
			name:     "Set/Fetch multiple coordinates",
			driverID: "1",
			coordinates: []Coordinates{
				{
					Lat:       1,
					Long:      2,
					UpdatedAt: common.Timestamp{Time: n.UTC()},
				},
				{
					Lat:       3,
					Long:      4,
					UpdatedAt: common.Timestamp{Time: n.UTC().Add(time.Minute)},
				},
			},
			expected: &[]Coordinates{
				{
					Lat:       1,
					Long:      2,
					UpdatedAt: common.Timestamp{Time: n.UTC()},
				},
				{
					Lat:       3,
					Long:      4,
					UpdatedAt: common.Timestamp{Time: n.UTC().Add(time.Minute)},
				},
			},
		},
		{
			name:        "Set/Fetch empty coordinates",
			driverID:    "1",
			coordinates: []Coordinates{},
			expected:    &[]Coordinates{},
		},
		{
			name:     "Set/Fetch multiple coordinates - ensure sorting order",
			driverID: "1",
			coordinates: []Coordinates{

				{
					Lat:       3,
					Long:      4,
					UpdatedAt: common.Timestamp{Time: n.UTC().Add(time.Minute)},
				},
				{
					Lat:       1,
					Long:      2,
					UpdatedAt: common.Timestamp{Time: n.UTC()},
				},
			},
			expected: &[]Coordinates{
				{
					Lat:       1,
					Long:      2,
					UpdatedAt: common.Timestamp{Time: n.UTC()},
				},
				{
					Lat:       3,
					Long:      4,
					UpdatedAt: common.Timestamp{Time: n.UTC().Add(time.Minute)},
				},
			},
		},
		{
			name:     "Set/Fetch multiple coordinates - ensure filtering last x minutes",
			driverID: "1",
			coordinates: []Coordinates{

				{
					Lat:       3,
					Long:      4,
					UpdatedAt: common.Timestamp{Time: n.UTC().Add(-10 * time.Minute)},
				},
				{
					Lat:       1,
					Long:      2,
					UpdatedAt: common.Timestamp{Time: n.UTC()},
				},
			},
			expected: &[]Coordinates{
				{
					Lat:       1,
					Long:      2,
					UpdatedAt: common.Timestamp{Time: n.UTC()},
				},
			},
		},
	}

	database := NewInMemoryDB(redis.NewClient(&redis.Options{}))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			for _, c := range test.coordinates {
				_ = database.Save(test.driverID, c, c.UpdatedAt.Time)
			}

			res, _ := database.Fetch(test.driverID, 5)

			if diff := deep.Equal(res, test.expected); diff != nil {
				t.Error(diff)
			}

			database.client.Del(test.driverID)
		})
	}
}
