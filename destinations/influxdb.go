package destinations

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Config struct {
	URL    string
	token  string
	org    string
	bucket string
}

type InfluxStore struct {
	Client influxdb2.Client
	Cfg    Config
	Writer api.WriteAPIBlocking
	Reader api.QueryAPI
}

func NewInfluxStore(ctx context.Context) InfluxStore {
	config := newInfluxCfg()
	client := influxdb2.NewClient(config.URL, config.token)

	store := InfluxStore{}

	store.setConfig(config)
	store.setClient(client)

	if !store.isOnline(ctx) {
		return InfluxStore{}
	}

	writer := client.WriteAPIBlocking(store.Cfg.org, store.Cfg.bucket)

	store.setWriter(writer)

	reader := client.QueryAPI(store.Cfg.org)
	store.setReader(reader)

	return store

}

func (s InfluxStore) isOnline(ctx context.Context) bool {
	online, err := s.Client.Ping(ctx)

	if online && err == nil {
		log.Println("Influx is online")
		return true
	}

	log.Printf("Influx is NOT online: %v", err)
	return false
}

func (s *InfluxStore) setConfig(c Config) {
	s.Cfg = c
}

func (s *InfluxStore) setClient(c influxdb2.Client) {
	s.Client = c
}

func (s *InfluxStore) setWriter(w api.WriteAPIBlocking) {
	s.Writer = w
}

func (s *InfluxStore) setReader(r api.QueryAPI) {
	s.Reader = r
}

func newInfluxCfg() Config {

	token := os.Getenv("INFLUX_TOKEN")
	url := os.Getenv("INFLUX_URL")
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_BUCKET")

	return Config{
		URL:    url,
		token:  token,
		org:    org,
		bucket: bucket,
	}
}

func (s InfluxStore) SendRecord(ctx context.Context, payload string) error {
	return nil
}
