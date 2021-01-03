package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/feitianlove/golib/common/logger"
	"github.com/sirupsen/logrus"
	"sync"
)

type Kafka struct {
	ServerAddr string
}

type KProduct struct {
	Producer sarama.SyncProducer
}
type KConsumer struct {
	Consumer sarama.Consumer
}

func NewKafkaProduct(kafka *Kafka) (*KProduct, error) {
	if kafka == nil {
		logger.Console.WithFields(logrus.Fields{
			"kafka": "kafka config is nil",
		}).Error("kafka config is nil")
		return nil, fmt.Errorf("the kafka struct is nil")
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true                   //是否开启消息发送成功后通知 successes channel
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机分区器
	client, err := sarama.NewClient([]string{kafka.ServerAddr}, config)
	if err != nil {
		logger.Console.WithFields(logrus.Fields{
			"kafka": fmt.Sprintf("%s", err),
		}).Error(" kafka NewClient error")
		return nil, err
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		logger.Console.WithFields(logrus.Fields{
			"kafka": fmt.Sprintf("%s", err),
		}).Error(" kafka NewSyncProducerFromClient error")
		return nil, err
	}
	return &KProduct{Producer: producer}, nil
}

func (k KProduct) SendMessage(data *sarama.ProducerMessage) error {
	produce := k.Producer
	partition, offset, err := produce.SendMessage(data)
	if err != nil {
		logger.Console.WithFields(logrus.Fields{
			"kafka": fmt.Sprintf("unable to produce message%s\n", err),
		}).Error(" kafka NewClient error")
		return err
	}
	logger.Console.WithFields(logrus.Fields{
		"data":      data,
		"partition": partition,
		"offset":    offset,
	}).Info("kafka SendMessage")
	return nil
}

func NewKafkaConsumer(kafka *Kafka) (*KConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{kafka.ServerAddr}, config)
	if err != nil {
		logger.Console.WithFields(logrus.Fields{
			"kafka": fmt.Sprintf("unable to NewKafkaConsumer%s\n", err),
		}).Error("NewKafkaConsumer")
		return nil, err
	}
	return &KConsumer{Consumer: consumer}, err
}

func (consumer KConsumer) RecvMessage(topic string) {
	var wg sync.WaitGroup
	c := consumer.Consumer
	partitionList, err := c.Partitions(topic)
	if err != nil {
		logger.Console.WithFields(logrus.Fields{
			"kafka": fmt.Sprintf("faild to get the list of partitions%s\n", err),
		}).Error("RecvMessage")
	}
	for partition := range partitionList {
		pc, err := c.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)
		if err != nil {
			logger.Console.WithFields(logrus.Fields{
				"kafka": fmt.Sprintf("Failed to start consumer for partition %d: %s\n", partition, err),
			}).Error("kakfa.go  RecvMessage error")
			return
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				logger.Console.WithFields(logrus.Fields{
					"Partition": msg.Partition,
					"Offset":    msg.Offset,
					"Key":       string(msg.Key),
					"Value":     string(msg.Value),
				}).Info("RecvMessage")
			}
		}(pc)
	}
	wg.Wait()
	defer func() {
		_ = c.Close()
	}()

}
