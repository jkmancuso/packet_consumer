package sources

import (
	"context"
	"crypto/tls"

	"fmt"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type kafkaConfig struct {
	topic     string
	partition int
	transport string
	host      string
	port      int
	tlsConfig *tls.Config
}

type kafkaConsumer struct {
	cfg    kafkaConfig
	reader *kafka.Reader
}

func NewKafkaConsumer(ctx context.Context, tlsConfigs ...*tls.Config) (kafkaConsumer, error) {
	kafkaCfg := newKafkaCfg(ctx)

	if len(tlsConfigs) != 0 {
		kafkaCfg.setTLS(tlsConfigs[0])
	}

	kConsumer := kafkaConsumer{
		cfg: kafkaCfg,
	}

	reader, err := kafkaCfg.getKafkaReader()

	if err != nil {
		log.Errorf("Err: %v\ncould not connect to kafka with params: %+v", err, kafkaCfg)
		return kConsumer, err
	}

	kConsumer.setReader(reader)

	log.Debugf("Returning kafka store: %+v", kConsumer)

	return kConsumer, nil
}

func (consumer kafkaConsumer) Teardown() {
	log.Printf("Tearing down kafka store")
	consumer.reader.Close()
}

func (consumer *kafkaConsumer) setReader(r *kafka.Reader) {
	consumer.reader = r
}

func (consumer kafkaConsumer) GetRecord(ctx context.Context) (string, error) {
	m, err := consumer.reader.ReadMessage(ctx)

	if err != nil {
		return "", err
	}

	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

	return string(m.Value), nil
}

func newKafkaCfg(_ context.Context) kafkaConfig {

	topic := os.Getenv("KAFKA_TOPIC")
	partition, _ := strconv.Atoi(os.Getenv("KAFKA_PARTITION"))
	transport := os.Getenv("KAFKA_TRANSPORT")
	host := os.Getenv("KAFKA_HOST")
	port, _ := strconv.Atoi(os.Getenv("KAFKA_PORT"))

	kafkaCfg := kafkaConfig{
		topic:     topic,
		partition: partition,
		transport: transport,
		host:      host,
		port:      port,
	}

	log.Debugf("Loading kafka cfg: %+v\n", kafkaCfg)

	return kafkaCfg

}

func (cfg *kafkaConfig) setTLS(tlsCfg *tls.Config) {
	log.Info("Setting tls config")
	cfg.tlsConfig = tlsCfg
}

// return a kafka connection handle
func (cfg *kafkaConfig) getKafkaReader() (*kafka.Reader, error) {

	log.Debugf("Connecting to kafka\n%+v", cfg)

	dialer := &kafka.Dialer{
		Timeout:   0,
		DualStack: true,
	}

	if cfg.tlsConfig != nil {
		dialer.TLS = cfg.tlsConfig
	}

	rc := kafka.ReaderConfig{
		Brokers:   []string{fmt.Sprintf("%s:%d", cfg.host, cfg.port)},
		Topic:     cfg.topic,
		Dialer:    dialer,
		Partition: cfg.partition,
	}

	err := rc.Validate()

	if err != nil {
		log.Error("Kafka read failed")
		return nil, err
	}

	r := kafka.NewReader(rc)

	log.Debugf("Success")

	return r, nil

}
