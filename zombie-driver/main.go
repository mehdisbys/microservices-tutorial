package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/heetch/MehdiSouilhed-technical-test/zombie-driver/app/domain"
	"github.com/heetch/MehdiSouilhed-technical-test/zombie-driver/app/handlers"
)

func main() {
	r := mux.NewRouter()

	c, err := domain.NewConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	// build driver service base url
	driverHost := fmt.Sprintf("%s/%s", c.DriverHost, c.DriverLocationsEndpoint)

	// fetcher is the client that makes the call to driver service
	fetcher := domain.NewLocationsFetcher(&http.Client{Timeout: c.Timeout * time.Second}, driverHost)

	// Defining how we will calculate distance - here it is Haversine method
	calculator := domain.HaversineDistance{}

	// Given a distance determine whether it is under the limit
	detector := domain.ZombieDetector{MinimumDistance: c.MinDistance}

	// RequestHandler is a struct holding the handler with its dependencies
	handler := handlers.NewRequestHandler(5, fetcher, calculator, detector)

	r.HandleFunc("/drivers/{id}", handler.GetDriverPings)

	log.Print("Listening on port 80")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 80), r))
}
