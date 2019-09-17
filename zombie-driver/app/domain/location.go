package domain

type LatLongerList []LatLonger

type LatLonger interface {
	GetLatitude() float64
	GetLongitude() float64
}

type CoordinatesList []Coordinates

type Coordinates struct {
	Lat       float64   `json:"latitude"` //Todo Validate
	Long      float64   `json:"longitude"`
	UpdatedAt Timestamp `json:"updated_at"`
}

func (c Coordinates) GetLatitude() float64 {
	return c.Lat
}

func (c Coordinates) GetLongitude() float64 {
	return c.Long
}
