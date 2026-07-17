package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	kafkaAddr := readEnvVar("KAFKA_BROKERS", "localhost:9093")
	httpPort := readEnvVar("PORT", "8082")

	kafka, err := NewKafkaManager(kafkaAddr)
	if err != nil {
		log.Printf("Не удалось инициализировать KafkaManager: %v", err)
		return
	}

	defer kafka.Close()

	for _, topic := range []string{TopicUsers, TopicMovies, TopicPayments} {
		err = kafka.ConsumeAllBackground(topic, consumeKafkaMessage)
		if err != nil {
			log.Printf("Не удалось запустить чтение топика '%v': %v", topic, err)
			return
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	runHttpRest(httpPort, kafka, stop)
	<-ctx.Done()
	shutdownHttpRest()
}

func readEnvVar(name string, fallback string) string {
	val := os.Getenv(name)
	if val == "" {
		val = fallback
		log.Printf("Не задана переменная окружения %v, используется значение %v", name, val)
	}
	return val
}

func consumeKafkaMessage(ev *Event) {
	if ev == nil {
		log.Printf("Получено пустое сообщение от Kafka")
		return
	}

	jsonData, err := json.Marshal(ev)
	if err != nil {
		log.Printf("Ошибка при обрабоке сообщения от Kafka: %v", err)
		return
	}

	fmt.Printf("Получено сообщение от Kafka: %v\n", string(jsonData))
}
