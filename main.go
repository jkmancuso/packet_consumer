package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := context.Background()

	setLogger()
	params := getCmdLineParams()
	_, err := NewSource(ctx, params["consumer"])

	if err != nil {
		log.Fatal("could not connect to the consumer!", err)
	}

}
