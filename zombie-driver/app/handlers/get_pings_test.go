package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-test/deep"
	"github.com/gorilla/mux"

	"github.com/heetch/MehdiSouilhed-technical-test/zombie-driver/app/domain"
	"github.com/heetch/MehdiSouilhed-technical-test/zombie-driver/app/domain/mocks"
)

func TestIsDriverZombie(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/drivers/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.

	tests := []struct {
		name         string
		coordinates  []domain.Coordinates
		zombie       bool
		returnErr    bool
		expectedCode int
	}{
		{
			name: "Is Zombie - immobile",
			coordinates: []domain.Coordinates{
				{Lat: 48.8566, Long: 2.3522},
				{Lat: 48.8566, Long: 2.3522},
				{Lat: 48.8566, Long: 2.3522},
			},
			zombie:       true,
			expectedCode: http.StatusOK,
		},
		{
			name: "Is not Zombie - more than 500",
			coordinates: []domain.Coordinates{
				{Lat: 48.8566, Long: 2.3522},
				{Lat: 51.5074, Long: 0.1278},
				{Lat: 35.6762, Long: 139.6503},
				{Lat: 41.3851, Long: 2.1734},
			},
			zombie:       false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Is Zombie - empty list",
			coordinates:  []domain.Coordinates{},
			zombie:       false,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Dependency returns error",
			returnErr:    true,
			expectedCode: http.StatusInternalServerError,
		},
	}

	calculator := domain.HaversineDistance{}
	detector := domain.ZombieDetector{MinimumDistance: 500}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			fetcher := mocks.NewMockFetcher(test.coordinates, test.returnErr)

			h := NewRequestHandler(5, fetcher, calculator, detector)

			handler := http.HandlerFunc(h.GetDriverPings)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			res := ZombieResponse{}

			_ = json.Unmarshal(rr.Body.Bytes(), &res)

			if diff := deep.Equal(res.Zombie, test.zombie); diff != nil {
				t.Error(diff)
			}
		})
	}
}
