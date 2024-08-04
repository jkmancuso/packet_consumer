package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/jkmancuso/packet_consumer/sources"
)

const kafkaEnvfile = "./kafka.env"

func NewSource(ctx context.Context, consumerType string) (sources.Consumer, error) {
	var source sources.Consumer
	var err error

	switch consumerType {
	case "kafka":
		loadEnv(kafkaEnvfile)

		if os.Getenv("KAFKA_TLS_ENABLED") == "TRUE" {

			log.Printf("Kafka is TLS enabled")

			tlsConfig, err := NewTLSConfig(
				os.Getenv("KAFKA_CLIENT_CERT"),
				os.Getenv("KAFKA_CLIENT_KEY"),
				os.Getenv("KAFKA_SERVER_CERT"))

			if err != nil {
				log.Fatalf("Unable to create TLS config: %v", err)
			}

			source, err = sources.NewKafkaConsumer(ctx, tlsConfig)

			if err != nil {
				log.Fatalf("Unable to get tls enabled %v store: %v", consumerType, err)
			}

		} else {
			log.Printf("Kafka is plaintext")

			source, err = sources.NewKafkaConsumer(ctx)

			if err != nil {
				log.Fatalf("Unable to get %v source: %v", consumerType, err)
			}

		}
	}

	return source, err

}
