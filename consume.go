package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/jkmancuso/packet_consumer/sources"
)

func startConsumer(ctx context.Context, consumer sources.Consumer) error {

	var err error
	var record string

	for {

		record, err = consumer.GetRecord(ctx)
		log.Printf("message %v\n", record)

		if err != nil {
			return err
		}
	}

}
