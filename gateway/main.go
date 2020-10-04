package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/heetch/MehdiSouilhed-technical-test/common"
	"github.com/heetch/MehdiSouilhed-technical-test/gateway/app/domain"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 5 * time.Second
	config.Producer.Return.Successes = true
	config.Admin.Timeout = time.Second * 30
	config.Admin.Retry.Max = 5
	config.Admin.Retry.Backoff = 5 * time.Second

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	time.Sleep(time.Second * 10)

	producer, err := sarama.NewSyncProducer([]string{"kafka1:9092"}, config)
	if err != nil {
		panic(err)
	}

	p := common.NewKafkaSender(producer)

	handler, err := domain.NewRequestHandler(p, &http.Client{Timeout: 5 * time.Second}, mux.NewRouter())
	if err != nil {
		panic(err)
	}

	configFile, err := domain.ParseFileConfig("config.yaml")

	if err != nil {
		log.Print(err)
		os.Exit(2)
	}

	handler.Gateway(configFile)
	log.Println("Listening on port 80")
	log.Println("Version 4")
	log.Fatal(http.ListenAndServe(":80", handler.GetRouter()))
}
