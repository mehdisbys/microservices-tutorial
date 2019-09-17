package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/heetch/MehdiSouilhed-technical-test/driver-location/app/domain"
	"github.com/heetch/MehdiSouilhed-technical-test/driver-location/app/handlers"
)

func main() {

	// Loading and validating configuration, if we dont have required config we will exit rather
	// than fail while the service is running
	c, err := domain.NewConfig("config.yaml")
	if err != nil {
		log.Print(os.Getwd())
		log.Print(err)
		os.Exit(1)
	}

	// Initialise database connection
	dbAddr := fmt.Sprintf("%s:%d", c.DatabaseHost, c.DatabasePort)
	log.Printf("Connecting to database at %s", dbAddr)
	database := domain.NewInMemoryDB(redis.NewClient(&redis.Options{Addr: dbAddr}))

	// Check the connection to DB
	err = database.Ping()
	if err != nil {
		log.Print(err)
		// having different exit code enables to localise errors quicker
		os.Exit(2)
	}

	// Instantiate queue handler
	s := handlers.NewSaveToDB(database)

	// Instantiate http router
	r := mux.NewRouter()

	// Register http handler
	handler := handlers.NewRequestHandler(database)
	r.HandleFunc("/drivers/{id}/locations", handler.GetDriverPings)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// ListenAndServe is a blocking operation - putting it in a goroutine
	go func() {
		defer wg.Done()
		port := 80
		fmt.Printf("Listening on port %d", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
	}()

	// Starting queue listener
	log.Printf("Queue listening on topic %s", c.QueueTopic)
	q, err := domain.NewNSQQueue(c.QueueTopic, "1", s)
	if err != nil {
		log.Printf("Could not create NSQ queue struct %s", err.Error())
		os.Exit(3)
	}

	err = q.Process(fmt.Sprintf("%s:%d", c.QueueHost, c.QueuePort))
	if err != nil {
		log.Printf("a queue error occurred : %s", err.Error())
		os.Exit(4)
	}
	wg.Wait()
}
