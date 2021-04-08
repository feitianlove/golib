package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
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
		return nil, fmt.Errorf("the kafka struct is nil")
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true                   //是否开启消息发送成功后通知 successes channel
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机分区器
	client, err := sarama.NewClient([]string{kafka.ServerAddr}, config)
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &KProduct{Producer: producer}, nil
}

func (k KProduct) SendMessage(data *sarama.ProducerMessage) (error, int32, int64) {
	produce := k.Producer
	partition, offset, err := produce.SendMessage(data)
	if err != nil {
		return err, 0, 0
	}
	return err, partition, offset
}

func NewKafkaConsumer(kafka *Kafka) (*KConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{kafka.ServerAddr}, config)
	if err != nil {
		return nil, err
	}
	return &KConsumer{Consumer: consumer}, err
}

func (consumer KConsumer) RecvMessage(returnData func(data *sarama.ConsumerMessage), topic string) error {
	var wg sync.WaitGroup
	c := consumer.Consumer
	partitionList, err := c.Partitions(topic)
	if err != nil {
		return err
	}
	for partition := range partitionList {
		pc, err := c.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)
		if err != nil {
			return err
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				returnData(msg)
				//logger.Console.WithFields(logrus.Fields{
				//	"Partition": msg.Partition,
				//	"Offset":    msg.Offset,
				//	"Key":       string(msg.Key),
				//	"Value":     string(msg.Value),
				//}).Info("RecvMessage")
			}
		}(pc)
	}
	wg.Wait()
	defer func() {
		_ = c.Close()
	}()
	return nil
}
