package domain

import (
	"time"

	"github.com/heetch/MehdiSouilhed-technical-test/common"
)

type CoordinatesList []Coordinates

type Coordinates struct {
	DriverID  string           `json:"courierId"`
	Lat       float64          `json:"latitude"` //Todo Validate
	Long      float64          `json:"longitude"`
	UpdatedAt common.Timestamp `json:"updated_at"`
}

func (c Coordinates) GetLatitude() float64 {
	return c.Lat
}

func (c Coordinates) GetLongitude() float64 {
	return c.Long
}

func (c *Coordinates) SetUpdatedAt(t time.Time) {
	c.UpdatedAt = common.Timestamp{Time: t}
}
