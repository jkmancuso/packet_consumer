package main

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	setLogger()
	params := getCmdLineParams()
	consumer, err := NewSource(ctx, params["consumer"])

	if err != nil {
		log.Fatal("could not connect to the consumer!", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		<-sigCh
		consumer.Teardown()
	}()

	if err = startConsumer(ctx, consumer); err != nil {
		log.Infof("finished consumer %v", err)
	}

}
