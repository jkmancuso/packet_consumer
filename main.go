package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/jkmancuso/packet_consumer/destinations"
	"github.com/jkmancuso/packet_consumer/sources"
	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	setLogger()
	params := getCmdLineParams()
	consumer, err := NewSource(ctx, *params["consumer"])

	if err != nil {
		log.Fatal("could not connect to the consumer!", err)
	}

	destination, err := NewDestination(ctx, *params["destination"])

	if err != nil {
		log.Fatal("could not connect to the destination!", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		<-sigCh
		consumer.Teardown()
	}()

	if err = Start(ctx, consumer, destination); err != nil {
		log.Infof("finished consumer %v", err)
	}

}

func Start(ctx context.Context, c sources.Consumer, d destinations.Destination) error {

	var err error
	var record string

	for {

		record, err = c.GetRecord(ctx)
		log.Printf("message %v\n", record)

		if err != nil {
			return err
		}
	}

}
