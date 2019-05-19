package main

import (
	"github.com/lt90s/goanalytics/conf"
	dataCodec "github.com/lt90s/goanalytics/event/codec"
	kafka2 "github.com/lt90s/goanalytics/event/pubsub/kafka"
	"github.com/lt90s/goanalytics/metric"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
)


func main() {
	subConfig := kafka.ReaderConfig{
		Topic:    conf.GetConfString(conf.TopicConfKey),
		Brokers:  conf.GetConfStringSlice(conf.KafkaBrokersConfKey),
		GroupID:  conf.GetConfString(conf.KafkaGroupIdConfKey),
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	}
	concurrentHandlerCount := conf.GetConfInt64(conf.KafkaNumberOfHandleProcessorConfKey)
	decoder := dataCodec.NewJsonCodec()
	subscriber := kafka2.NewSubscriber(subConfig, decoder, concurrentHandlerCount)
	subscriber.Start()
	defer subscriber.Shutdown()

	metric.SetupMetricProcessor(subscriber)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	<-sigChan
}
