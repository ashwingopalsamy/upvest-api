package main

import (
	"github.com/ashwingopalsamy/upvest-api/internal/event"
	log "github.com/sirupsen/logrus"
)

var publisher *event.Publisher

func initKafkaPublisher() {
	publisher = event.NewPublisher("kafka:9092", "user-events")
	log.Info("Kafka publisher initialized")
}
