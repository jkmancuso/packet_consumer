package destinations

import (
	"context"
	"os"
	"strings"

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

	if err != nil {
		if strings.Contains(err.Error(), "tls: failed to verify certificate") {
			log.Println("Influx is online but with self signed crt")
			return true
		}

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

	url := os.Getenv("INFLUX_URL")
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_BUCKET")

	//this needs to be in your ENV variables either in your
	//OS, container, or supporting application
	token := os.Getenv("INFLUX_TOKEN")

	if len(url) == 0 ||
		len(org) == 0 ||
		len(bucket) == 0 ||
		len(token) == 0 {
		log.Errorf("!Missing influx variable!")
	}

	cfg := Config{
		URL:    url,
		token:  token,
		org:    org,
		bucket: bucket,
	}

	log.Debugf("Influx Config: %+v", cfg)

	return cfg
}

func (s InfluxStore) SendRecord(ctx context.Context, payload string) error {
	return nil
}
