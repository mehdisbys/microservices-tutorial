package domain

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heetch/MehdiSouilhed-technical-test/common"
	"github.com/rs/zerolog/log"
)

type QueueProducer interface {
	Publish(topic string, body []byte) error
}

type RequestHandler struct {
	client   *http.Client
	producer common.Sender
	router   *mux.Router
}

type Message struct {
	Body       []byte            `json:"body"`
	Parameters map[string]string `json:"parameters"`
}

const (
	logTraceID = "traceID"
)

func NewRequestHandler(p common.Sender, client *http.Client, r *mux.Router) (*RequestHandler, error) {

	return &RequestHandler{
		producer: p,
		client:   client,
		router:   r,
	}, nil
}

func (s *RequestHandler) GetRouter() *mux.Router {
	return s.router
}

func (s *RequestHandler) Gateway(config Config) {

	for _, c := range config.Urls {
		switch {
		case c.Nsq != nil:
			s.makeAsyncHandler(c.Method, c.Path, c.Nsq.Topic)

		case c.HTTP != nil:
			host := c.HTTP.Host
			s.makeSyncHandler(c.Method, c.Path, host)
		}
	}
}

func (s *RequestHandler) makeAsyncHandler(method, path, topic string) {

	log.Info().Msgf("Registering async handler for [method|path|topic]: [%s|%s|%s]", method, path, topic)

	s.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		if r.Method != method {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		traceID := common.ExtractTraceIDFromReq(r)

		request, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID)
		}

		urlVars := mux.Vars(r)
		// Pass the traceID downstream
		urlVars[common.TraceIDHeader] = traceID

		m := Message{Body: request, Parameters: urlVars}

		mbytes, err := json.Marshal(m)
		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID)
		}

		log.Info().Interface("params", urlVars).Str("topic", topic).Msg("transforming request to async event")
		err = s.producer.Send(topic, string(mbytes))
		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID).
				Msg("could not publish message")

		}
		w.WriteHeader(http.StatusOK)
	})
}

func (s *RequestHandler) makeSyncHandler(method, path, host string) {

	log.Info().Msgf("Registering http proxy handler for [method|path|host]: [%s|%s|%s]", method, path, host)

	s.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		if r.Method != method {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		traceID := common.ExtractTraceIDFromReq(r)

		res, err := s.proxy("http://"+host+r.URL.Path, r)
		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer res.Body.Close()
		respBytes, err := ioutil.ReadAll(res.Body)

		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(res.StatusCode)

		_, err = w.Write(respBytes)
		if err != nil {
			log.Error().Err(err).Str(logTraceID, traceID)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func (s *RequestHandler) proxy(proxyURL string, r *http.Request) (*http.Response, error) {

	req, err := http.NewRequest(r.Method, proxyURL, r.Body)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return nil, err
	}

	params := r.URL.Query()
	req.URL.RawQuery = params.Encode()

	// Pass the traceID downstream
	req.Header.Add(common.TraceIDHeader, common.ExtractTraceIDFromReq(r))

	response, err := s.client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("error")
		return nil, err
	}

	return response, nil
}
