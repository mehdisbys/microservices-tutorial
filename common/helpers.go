package common

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

const TraceIDHeader = "X-Trace-Id"

func GetIntVariableValue(r *http.Request, key string) (int, error) {
	vars := mux.Vars(r)
	strValue, ok := vars[key]
	if !ok || strValue == "" {
		return 0, errors.New("missing value")
	}

	intValue, err := strconv.Atoi(strValue)

	if err != nil {
		return 0, err
	}

	return intValue, nil
}

func GetIntParamValue(r *http.Request, key string) (int, error) {
	strValue := r.URL.Query().Get(key)

	if strValue == "" {
		return 0, nil
	}

	intValue, err := strconv.Atoi(strValue)

	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func ExtractTraceIDFromReq(r *http.Request) (traceID string) {
	traceID = r.Header.Get(TraceIDHeader)
	if traceID == "" {
		traceID = uuid.NewV4().String()
	}
	return traceID
}
