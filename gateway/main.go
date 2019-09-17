package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/nsqio/go-nsq"

	"github.com/heetch/MehdiSouilhed-technical-test/gateway/app/domain"
)

func main() {
	addr := fmt.Sprintf("nsqd:%d", 4150)
	config := nsq.NewConfig()

	p, err := nsq.NewProducer(addr, config)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

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
