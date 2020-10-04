package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/heetch/MehdiSouilhed-technical-test/common"
	"github.com/heetch/MehdiSouilhed-technical-test/driver-location/app/domain"
	"github.com/heetch/MehdiSouilhed-technical-test/driver-location/app/handlers"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	time.Sleep(10 * time.Second)

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
	topic := "locations"
	//	partition := 0
	//	offsetType := kingpin.Flag("offsetType", "Offset Type (OffsetNewest | OffsetOldest)").Default("-1").Int()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"kafka1:9092"}, config)
	if err != nil {
		log.Panic(err)
	}

	stream := common.NewKafkaConsumer(consumer, s)

	stream.Receive(topic)

	wg.Wait()
}
