package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/jkmancuso/packet_consumer/destinations"
	"github.com/jkmancuso/packet_consumer/sources"
	log "github.com/sirupsen/logrus"
)

const measurement = "ip"

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
	var record []byte
	var prevTime time.Time
	aggChan := make(chan string, 1)

	for {

		record, err = c.GetRecord(ctx)
		log.Printf("message %v\n", string(record))

		if err != nil {
			log.Error(err)
			return err
		}

		ipEntry := ipEntry{}

		if err = json.Unmarshal(record, &ipEntry); err != nil {
			log.Error(err)
			return err
		}

		//run every minute (except the first iteration)
		if !prevTime.IsZero() && ipEntry.getTime().Minute() != prevTime.Minute() {
			go d.Aggregate(ctx, prevTime, ipEntry.getTime(), aggChan)
		}

		err = d.SendRecord(ctx,
			measurement,
			ipEntry.getTags(),
			ipEntry.getFields(),
			ipEntry.getTime())

		if err != nil {
			log.Error(err)
			return err
		}

		prevTime = ipEntry.getTime()
	}

}
