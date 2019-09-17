package mocks

import (
	"errors"

	"github.com/heetch/MehdiSouilhed-technical-test/zombie-driver/app/domain"
)

type MockFetcher struct {
	coordinates []domain.Coordinates
	returnError bool
}

func NewMockFetcher(c []domain.Coordinates, returnError bool) domain.Fetcher {
	return MockFetcher{
		coordinates: c,
		returnError: returnError,
	}
}

func (m MockFetcher) GetList(courierID int, minutes int) ([]domain.Coordinates, error) {
	if m.returnError {
		return nil, errors.New("failed to get pings")
	}

	return m.coordinates, nil
}
