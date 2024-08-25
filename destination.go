package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/jkmancuso/packet_consumer/destinations"
)

const influxdbEnvfile = "./influxdb.env"

// return Interface
func NewDestination(ctx context.Context, storeType string) (destinations.Destination, error) {
	var store destinations.Destination
	var err error

	log.Printf("Using store type %v", storeType)

	switch storeType {
	case "influxdb":
		loadEnv(influxdbEnvfile)

		//return concrete type
		store = destinations.NewInfluxStore(ctx)

	}

	return store, err

}
