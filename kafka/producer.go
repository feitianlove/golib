package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func SendMessageToKafka() {
	config := Kafka{
		//ServerAddr: []string{"9.135.147.57:9093", "9.135.147.57:9092", "9.135.147.57:9094"},
		ServerAddr: []string{"9.135.147.57:9092"},

		ProducerRetryParam: &ProducerRetry{
			Max:         3,
			Backoff:     200 * time.Millisecond,
			BackoffFunc: nil,
		},
		ProducerAck: -1,
	}
	client, err := NewKafkaClient(&config, Producer, nil)
	if err != nil {
		panic(err)
	}
	// 这里注意如果topic不存在默认指定创建 PartitionCount: 1 , ReplicationFactor: 1的topic
	// 如果topic PartitionCount: 3,ReplicationFactor: 1 如果其他两个partition 挂了，kafka只会往其中一个存活的kafka分区中写数据。
	for i := 1; i < 100; i++ {
		time.Sleep(time.Second * 1)
		err, partition, offset := client.SendMessage(&sarama.ProducerMessage{
			Topic:     "ftfeng-producer-test-5",
			Key:       nil,
			Value:     sarama.StringEncoder(fmt.Sprintf("ftfeng-producer-value-%d", i)),
			Headers:   nil,
			Metadata:  nil,
			Offset:    0,
			Partition: 0,
			Timestamp: time.Time{},
		})
		if err != nil {
			fmt.Printf("ftfeng producer value %d faild, err=%s\n", i, err.Error())
		} else {
			fmt.Printf("ftfeng producer value %d success,partition=%d, offset=%d\n", i, partition, offset)
		}
	}

}
