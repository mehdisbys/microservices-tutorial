package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/heetch/MehdiSouilhed-technical-test/common"
	"github.com/heetch/MehdiSouilhed-technical-test/zombie-driver/app/domain"
	"github.com/rs/zerolog/log"
)

type RequestHandler struct {
	intervalMinutes int
	fetcher         domain.Fetcher
	calculator      domain.DistanceEstimator
	detector        domain.ZombieDetector
}

type ZombieResponse struct {
	ID     int  `json:"id"`
	Zombie bool `json:"zombie"`
}

const (
	courierID      = `id`
	minutes        = `minutes`
	defaultMinutes = 5
	logTraceID     = "traceID"
)

func NewRequestHandler(
	interval int,
	fetcher domain.Fetcher,
	estimator domain.DistanceEstimator,
	detector domain.ZombieDetector) RequestHandler {
	return RequestHandler{
		intervalMinutes: interval,
		fetcher:         fetcher,
		calculator:      estimator,
		detector:        detector,
	}
}

// GetDriversPings will fetch pings from the relevant service
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

	pings, err := s.fetcher.GetList(id, m)

	if err != nil {
		log.Error().Err(err).Str(logTraceID, traceID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := ZombieResponse{ID: id, Zombie: false}

	// If there are not a least two coordinates we cannot tell if driver is a zombie
	if len(pings) <= 1 {
		log.Info().Str(logTraceID, traceID).Msgf("driver %d does not have enough pings to check if he is zombie (count=%d)", id, len(pings))
		res, err := json.Marshal(response)
		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = w.Write(res)
		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID)
		}
		return
	}

	// Otherwise get the distance that the driver has done
	distance := s.calculator.SumDistanceList(pings)
	isZombie := s.detector.IsZombie(distance)
	response.Zombie = isZombie

	log.Info().Str(logTraceID, traceID).Msgf("driver %d has gone %.2f meters the last %d minutes "+
		"when minimum is %.2f meters and therefore zombie = %t",
		id, distance, m, s.detector.MinimumDistance, isZombie)

	res, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Str(logTraceID, traceID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(res)

	if err != nil {
		log.Error().Err(err).Str(logTraceID, traceID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
