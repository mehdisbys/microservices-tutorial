package handlers

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/heetch/MehdiSouilhed-technical-test/driver-location/app/domain"
	"github.com/nsqio/go-nsq"
	. "github.com/onsi/gomega"
)

type MockDB struct {
	store map[string][]domain.Coordinates
}

func (m *MockDB) Save(driverID string, coordinates domain.Coordinates, time time.Time) error {
	coordinates.SetUpdatedAt(time)
	m.store[driverID] = append(m.store[driverID], coordinates)
	return nil
}

func (m MockDB) Fetch(driverID string, minutes int) (*[]domain.Coordinates, error) {
	val := m.store[driverID]
	return &val, nil
}

func (m MockDB) Ping() error {
	return nil
}

func TestHandleMessage(t *testing.T) {

	coord, _ := json.Marshal(domain.Coordinates{Lat: 1, Long: 2})
	m := &MockDB{store: map[string][]domain.Coordinates{}}
	handler := NewSaveToDB(m)
	g := NewGomegaWithT(t)

	tests := []struct {
		name      string
		driverID  string
		message   domain.Message
		expected  []domain.Coordinates
		expectErr interface{}
	}{
		{
			name:     "valid message gets stored to db",
			driverID: "10",
			message: domain.Message{
				Body: coord,
				Parameters: map[string]string{
					"id": "10", // driver ID
				},
			},
			expected: []domain.Coordinates{
				{
					DriverID: "10",
					Lat:      1,
					Long:     2,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			body, err := json.Marshal(test.message)
			if err != nil && test.expectErr == nil {
				t.Log(err)
			}

			// HandleMessage will save it to db
			err = handler.HandleMessage(&nsq.Message{Body: body})
			if err != nil && test.expectErr == nil {
				t.Log(err)
			}

			// Fetching it from the mock DB
			res, err := m.Fetch(test.driverID, 0)
			if err != nil {
				t.Log(err)
			}

			for i, r := range *res {
				g.Expect(r.UpdatedAt.Time).Should(BeTemporally("~", time.Now()))
				g.Expect(r.Lat).To(Equal(test.expected[i].Lat))
				g.Expect(r.Long).To(Equal(test.expected[i].Long))
				g.Expect(r.DriverID).To(Equal(test.expected[i].DriverID))
			}
		})
	}
}

func TestHandleMessageErrors(t *testing.T) {
	coord, _ := json.Marshal(domain.Coordinates{Lat: 1, Long: 2})
	m := &MockDB{store: map[string][]domain.Coordinates{}}
	handler := NewSaveToDB(m)

	tests := []struct {
		name      string
		driverID  string
		message   interface{}
		expected  []domain.Coordinates
		expectErr interface{}
	}{
		{
			name:     "valid message - missing driver id",
			driverID: "10",
			message: domain.Message{
				Body: coord,
			},
			expectErr: MissingDriverID{},
		},
		{
			name:      "invalid message type",
			driverID:  "10",
			message:   []string{""},
			expectErr: &json.UnmarshalTypeError{},
		},
		{
			name:     "invalid body inside message",
			driverID: "10",
			message: domain.Message{
				Body: []byte{},
			},
			expectErr: &json.SyntaxError{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			body, err := json.Marshal(test.message)
			if err != nil && test.expectErr == nil {
				t.Log(err)
			}

			// HandleMessage will save it to db
			err = handler.HandleMessage(&nsq.Message{Body: body})
			if err != nil && test.expectErr == nil {
				t.Log(err)
			}

			if test.expectErr != nil && err == nil {
				t.Errorf("was expecting error of type %s but got none", reflect.TypeOf(test.expectErr).String())
			}

			if test.expectErr != nil && err != nil {
				expectedType := reflect.TypeOf(test.expectErr).String()
				actualType := reflect.TypeOf(err).String()

				if actualType != expectedType {
					t.Errorf("was expecting error of type %s but got error type %s", expectedType, actualType)
				}
				return
			}
		})
	}
}
