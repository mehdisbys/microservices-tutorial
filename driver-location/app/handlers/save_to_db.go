package handlers

import (
	"encoding/json"
	"time"

	"github.com/heetch/MehdiSouilhed-technical-test/common"
	"github.com/heetch/MehdiSouilhed-technical-test/driver-location/app/domain"
	"github.com/nsqio/go-nsq"
	"github.com/rs/zerolog/log"
)

// SaveToDB holds the dependencies for the queue handler
// it implements the nsq.Handler interface
type SaveToDB struct {
	database domain.DB
}

func NewSaveToDB(db domain.DB) *SaveToDB {
	return &SaveToDB{
		database: db,
	}
}

// MissingDriverID is a custom error type returned when a queue message is missing the
// driverID
type MissingDriverID struct {
	message string
}

func (m MissingDriverID) Error() string {
	return m.message
}

// HandleMessage will unmarshal a message from the queue and save it to the database
func (s SaveToDB) HandleMessage(message *nsq.Message) error {
	if message != nil {
		m := domain.Message{}

		// First unmarshal the envelope
		err := json.Unmarshal(message.Body, &m)
		if err != nil {
			log.Error().Err(err)
			return err
		}

		traceID := m.Parameters[common.TraceIDHeader]
		location := domain.Coordinates{}

		// Second unmarshal the content of the message
		err = json.Unmarshal(m.Body, &location)
		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID)
			return err
		}

		if id, ok := m.Parameters["id"]; ok {
			location.DriverID = id
		} else {
			return MissingDriverID{"no driver id found in message"}
		}

		log.Info().
			Interface("message", location).
			Str("driverID", location.DriverID).
			Str(logTraceID, traceID).
			Msg("saving location to db")

		return s.database.Save(location.DriverID, location, time.Now())
	}
	return nil
}
