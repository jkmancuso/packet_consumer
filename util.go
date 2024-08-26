package main

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func loadEnv(paths ...string) {
	if len(paths) == 0 {
		godotenv.Load()
	} else {
		for _, path := range paths {
			godotenv.Load(path)
		}
	}

	if _, isCI := os.LookupEnv("GITHUB_ACTIONS"); isCI {
		os.Setenv("KAFKA_TLS_ENABLED", "FALSE")
	}

}

func setLogger() {
	loadEnv()

	level, err := log.ParseLevel(os.Getenv("LEVEL"))

	if err != nil {
		log.Panic("Unable to recognize logging")
	}

	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	//log.SetReportCaller(true)

	log.Printf("log level set to: %v", level.String())

}

func getCmdLineParams() map[string]*string {
	params := make(map[string]*string)

	params["consumer"] = flag.String("consumer", "kafka", "")
	params["destination"] = flag.String("destination", "influxdb", "")

	flag.Parse()

	return params
}
