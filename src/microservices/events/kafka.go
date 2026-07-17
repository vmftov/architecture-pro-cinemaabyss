package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type KafkaManager struct {
	BrokerAddrs          []string
	Config               *sarama.Config
	Producer             sarama.SyncProducer
	ConsumerGroups       []sarama.ConsumerGroup
	ConsumerGroupsMutex  *sync.Mutex
	ConsumerGroupsWg     *sync.WaitGroup
	ConsumerGroupsCtx    context.Context
	ConsumerGroupsCancel func()
}

func (m *KafkaManager) Close() {
	if m.ConsumerGroupsCancel != nil {
		m.ConsumerGroupsCancel()
	}

	if m.ConsumerGroupsWg != nil {
		m.ConsumerGroupsWg.Wait()
	}

	m.ConsumerGroupsMutex.Lock()
	if m.ConsumerGroups != nil {
		for _, cg := range m.ConsumerGroups {
			err := cg.Close()
			if err != nil {
				log.Printf("Ошибка при закрытии ConsumerGroup: %v", err)
			}
		}
	}
	m.ConsumerGroupsMutex.Unlock()

	if m.Producer != nil {
		err := m.Producer.Close()
		if err != nil {
			log.Printf("Ошибка при закрытии Producer: %v", err)
		}
	}
}

type ConsumeGroupHandler struct {
	Callback func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

func NewKafkaManager(addr string) (*KafkaManager, error) {
	checkerr := awaitKafka(addr)
	if checkerr != nil {
		return nil, checkerr
	}

	m := &KafkaManager{}

	m.BrokerAddrs = []string{addr}
	m.Config = createKafkaConfig()
	m.ConsumerGroups = make([]sarama.ConsumerGroup, 0)
	m.ConsumerGroupsMutex = &sync.Mutex{}
	m.ConsumerGroupsWg = &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	m.ConsumerGroupsCtx = ctx
	m.ConsumerGroupsCancel = cancel

	producer, err := sarama.NewSyncProducer(m.BrokerAddrs, m.Config)
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания Kafka Producer: %w", err)
	}
	m.Producer = producer

	return m, nil
}

func awaitKafka(addr string) error {
	addrs := []string{addr}
	config := createKafkaConfig()

	testKafkaConnection := func() bool {
		admin, err := sarama.NewClusterAdmin(addrs, config)
		if err != nil {
			return false
		}

		defer admin.Close()

		_, err = admin.ListTopics()
		return err == nil
	}

	retries := 30
	delay := 2 * time.Second

	log.Printf("Проверка готовности Kafka")

	for i := 0; i < retries; i++ {
		if testKafkaConnection() {
			log.Printf("Подключение к Kafka выполнено")
			return nil
		}

		log.Printf("Не удалось подключиться к Kafka. Повторная попытка через 2 сек")
		time.Sleep(delay)
	}

	return fmt.Errorf("Не удалось подключиться к Kafka")
}

func createKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Metadata.AllowAutoTopicCreation = true
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true
	return config
}

func (m *KafkaManager) Produce(topic string, ev Event) (int32, int64, error) {
	jsonData, err := json.Marshal(ev)
	if err != nil {
		return 0, 0, err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(jsonData),
	}

	return m.Producer.SendMessage(msg)
}

func (ConsumeGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	log.Println("Consumer group setup")
	return nil
}

func (ConsumeGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer group cleanup")
	return nil
}

func (h ConsumeGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return h.Callback(sess, claim)
}

func (m *KafkaManager) ConsumeAllBackground(topic string, consumeEvent func(*Event)) error {
	cg, cgerr := sarama.NewConsumerGroup(m.BrokerAddrs, "event-service."+topic+".consumer", m.Config)
	if cgerr != nil {
		return fmt.Errorf("Ошибка создания Kafka ConsumerGroup: %w", cgerr)
	}

	m.ConsumerGroupsMutex.Lock()
	m.ConsumerGroups = append(m.ConsumerGroups, cg)
	m.ConsumerGroupsMutex.Unlock()

	ctx := m.ConsumerGroupsCtx

	// чтение ошибок consumer group

	m.ConsumerGroupsWg.Add(1)

	go func() {
		defer m.ConsumerGroupsWg.Done()

		for {
			select {
			case err, ok := <-cg.Errors():
				if !ok {
					return
				}
				if err != nil {
					log.Printf("Ошибка в Kafka ConsumerGroup (топик '%v'): %v", topic, err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// чтение сообщений

	m.ConsumerGroupsWg.Add(1)

	go func() {
		defer m.ConsumerGroupsWg.Done()

		handler := ConsumeGroupHandler{
			Callback: func(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
				for msg := range claim.Messages() {
					processMessage(sess, msg, topic, consumeEvent)
				}
				return nil
			},
		}

		for {
			err := cg.Consume(ctx, []string{topic}, handler)
			if ctx.Err() != nil {
				log.Printf("Kafka consumer loop остановлен: %v", ctx.Err())
				break
			}
			if err != nil {
				log.Printf("Ошибка в ConsumeGroup: %v", err)
				select {
				case <-time.After(time.Second):
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return nil
}

func processMessage(sess sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage, topic string, consumeEvent func(*Event)) {
	defer sess.MarkMessage(msg, "")

	defer func() {
		r := recover()
		if r != nil {
			log.Printf("Ошибка (panic) при обработке сообщения из топика '%v': %v", topic, r)
		}
	}()

	var ev *Event
	err := json.Unmarshal(msg.Value, &ev)
	if err != nil {
		log.Printf("Ошибка при парсинге сообщения: %v", err)
		return
	}
	if ev == nil {
		log.Printf("Пропущено пустое сообщение (null) из топика '%s'", topic)
		return
	}

	consumeEvent(ev)
}
