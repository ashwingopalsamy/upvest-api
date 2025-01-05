package main

import (
	"github.com/ashwingopalsamy/upvest-api/internal/kafka"
	log "github.com/sirupsen/logrus"
)

var publisher *kafka.Publisher

func initKafkaPublisher() {
	publisher = kafka.NewPublisher("kafka:9092", "user-events")
	log.Info("Kafka publisher initialized")
}
