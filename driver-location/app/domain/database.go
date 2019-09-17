package domain

import (
	"encoding/json"
	"log"
	"sort"
	"time"

	"github.com/go-redis/redis"
)

// DB is an interface to a database where pings will be stored
type DB interface {
	Save(driverID string, coordinates Coordinates, time time.Time) error
	Fetch(driverID string, minutes int) (*[]Coordinates, error)
	Ping() error
}

type InMemoryDB struct {
	client *redis.Client
}

func NewInMemoryDB(client *redis.Client) *InMemoryDB {
	return &InMemoryDB{
		client: client,
	}
}

func (d *InMemoryDB) Ping() error {
	_, err := d.client.Ping().Result()
	return err
}

// Save takes an updatedAt value and persists coordinates for a driverID
func (d *InMemoryDB) Save(driverID string, coordinates Coordinates, time time.Time) error {

	coordinates.SetUpdatedAt(time)
	c, err := json.Marshal(coordinates)
	if err != nil {
		return err
	}

	log.Printf("saving coordinates for driver %s", driverID)

	return d.client.SAdd(driverID, c).Err()
}

// Fetch retrieves coordinates for a driverID given they are not older than `minutes`
func (d *InMemoryDB) Fetch(driverID string, minutes int) (*[]Coordinates, error) {

	res, err := d.client.SMembers(driverID).Result()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	coords := []Coordinates{}

	for _, r := range res {

		c := Coordinates{}

		if err := json.Unmarshal([]byte(r), &c); err != nil {
			return nil, err
		}

		diff := now.Sub(c.UpdatedAt.Time)

		// Only take the last x minutes
		if diff < (time.Minute * time.Duration(minutes)) {
			coords = append(coords, c)
		}
	}

	sort.Slice(coords, func(i, j int) bool {
		return coords[i].UpdatedAt.Before(coords[j].UpdatedAt.Time)
	})

	return &coords, nil
}
