package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func ReceiveMessageFromKafka(topic string) {
	config := Kafka{
		ServerAddr: []string{"9.135.147.57:9093", "9.135.147.57:9092", "9.135.147.57:9094"},
		ConsumerOffsetParam: &ConsumerOffsetParam{
			// 这里虽然设置了禁止自动提交，但是应该还是自动提交了，没用，如果需要设置禁用自动提交还是需要使用consumerGroup
			AutoCommit:            false,
			RetryBackoff:          0,
			RetryBackoffFunc:      nil,
			ConsumerDefaultOffset: sarama.OffsetOldest,
			MAX:                   3,
		},
	}
	client, err := NewKafkaClient(&config, Consumer, nil)
	if err != nil {
		panic(err)
	}
	err = client.ReceiveMessageByConsumer(topic, func(data *sarama.ConsumerMessage, ch <-chan *sarama.ConsumerError) {
		time.Sleep(time.Second * 2)
		go func() {
			err := <-ch
			fmt.Printf("error from kafka consumer: %s\n", err)

		}()
		fmt.Printf("ftfeng consumer topic=%s, partion=%d, offset=%d, value=%s\n", topic, data.Partition, data.Offset, data.Value)
	})
	if err != nil {
		panic(err)
	}
}

func ReceiveMessageFromKafkaByConsumerGroup() {
	config := Kafka{
		ServerAddr: []string{"9.135.147.57:9093", "9.135.147.57:9092", "9.135.147.57:9094"},
		ConsumerOffsetParam: &ConsumerOffsetParam{
			// 这里虽然设置了禁止自动提交，但是应该还是自动提交了，没用，如果需要设置禁用自动提交还是需要使用consumerGroup
			AutoCommit:            false,
			RetryBackoff:          0,
			RetryBackoffFunc:      nil,
			ConsumerDefaultOffset: sarama.OffsetOldest,
			MAX:                   3,
		},
	}
	client, err := NewKafkaClient(&config, Consumer, nil)
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	var cc = ConsumerByGroup{}
	err = client.ReceiveMessageByConsumerGroup("ftfeng-last-test", ctx, []string{"ftfeng-producer-test-5"}, cc)
	if err != nil {
		panic(err)
	}
}
