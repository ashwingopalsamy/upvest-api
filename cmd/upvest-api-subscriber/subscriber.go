package main

import (
	"encoding/json"
	"fmt"

	"github.com/ashwingopalsamy/upvest-api/internal/kafka"
	log "github.com/sirupsen/logrus"
)

var subscriber *kafka.Subscriber

func initKafkaSubscriber() {
	subscriber = kafka.NewSubscriber("kafka:9092", "user-events", "user-subscriber-group")
	log.Info("Kafka subscriber initialized")
}

func kafkaListener(key, value []byte) error {
	var event map[string]interface{}
	if err := json.Unmarshal(value, &event); err != nil {
		log.Errorf("failed to unmarshal message: %v", err)
		return err
	}

	action := event["action"]
	log.Infof("Processing event: %v", action)

	fmt.Printf("Processed event: %v\n", event)
	return nil
}
