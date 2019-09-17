package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/heetch/MehdiSouilhed-technical-test/common"
	"github.com/heetch/MehdiSouilhed-technical-test/driver-location/app/domain"
	"github.com/rs/zerolog/log"
)

type RequestHandler struct {
	database domain.DB
}

const (
	courierID      = `id`
	minutes        = `minutes`
	defaultMinutes = 5
)

func NewRequestHandler(database domain.DB) *RequestHandler {
	return &RequestHandler{
		database: database,
	}
}

const (
	logTraceID = "traceID"
)

// GetDriversPings will fetch pings from the database
func (s *RequestHandler) GetDriverPings(w http.ResponseWriter, r *http.Request) {
	traceID := common.ExtractTraceIDFromReq(r)

	id, err := common.GetIntVariableValue(r, courierID)
	if err != nil {
		log.Error().Err(err).Str(logTraceID, traceID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := common.GetIntParamValue(r, minutes)
	if err != nil {
		log.Error().Err(err).Str(logTraceID, traceID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if m == 0 {
		m = defaultMinutes
	}

	pings, err := s.database.Fetch(strconv.Itoa(id), m)

	if err != nil {
		log.Error().Err(err).Str(logTraceID, traceID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(pings)
	if err != nil {
		log.Error().Err(err).Str(logTraceID, traceID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(response)
	if err != nil {
		log.Error().Err(err).Str(logTraceID, traceID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
