package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"golang.org/x/net/context"
	"sync"
	"time"
)

//Producer 生产消息的retry的参数
type ProducerRetry struct {
	Max     int
	Backoff time.Duration
	// 这个参数用于 更加复杂的重试策略，如果设置优先级高于Backoff参数
	BackoffFunc func(retries, maxRetries int) time.Duration
}

/*
	Retention 参数默认被禁用， 禁用的情况下使用offsets.retention.minutes
	offsets.retention.minutes则是记录topic的偏移量日志的保留时长。
	偏移量是指向消费者已消耗的最新消息的指针。 比如，你消费了10条消息，那么偏移量将移动10个位置。 这个偏移量会被记录到日志中，以便我们下次消费时
	知道应该从哪个offset开始继续消费。
*/

type Kafka struct {
	ServerAddr          []string
	ProducerRetryParam  *ProducerRetry
	ProducerAck         sarama.RequiredAcks
	ConsumerOffsetParam *ConsumerOffsetParam
}

// Consumer 生产环境retry参数
type ConsumerOffsetParam struct {
	// 是否开启自动提交,仅对consumerGroup生效
	AutoCommit bool
	// 自动提交的默认间隔，仅在Enable = TRUE生效
	AutoCommitInterval time.Duration
	//重试的间隔参数
	RetryBackoff time.Duration
	// 这个参数用于 更加复杂的重试策略，如果设置优先级高于Backoff参数
	RetryBackoffFunc func(retries int) time.Duration
	// 消费者默人的偏移量， 如果之前没有提交过默认为 OffsetNewest， 只能为默认的两个值。
	ConsumerDefaultOffset int64
	// The total number of times to retry failing commit
	MAX int
}

type ClientKafka struct {
	ConsumerClient      sarama.Client
	ConsumerGroupClient sarama.Client
	ProducerClient      sarama.Client
}

type RoleType string

const (
	Consumer     RoleType = "CONSUMER"
	Producer     RoleType = "PRODUCER"
	ConsumeGroup RoleType = "CONSUME_GROUP"
)

type setCustomParamFunc func(config *sarama.Config)

func NewKafkaClient(kafka *Kafka, roleType RoleType, setCustomParam setCustomParamFunc) (*ClientKafka, error) {
	switch roleType {
	case Consumer:
		return newKafkaConsumerClient(kafka, setCustomParam)
	case Producer:
		return newKafkaProductClient(kafka, setCustomParam)
	case ConsumeGroup:
		return newKafkaConsumerClient(kafka, setCustomParam)
	default:
		return nil, fmt.Errorf("don't support this RoteType %s", roleType)
	}
}

func newKafkaProductClient(kafka *Kafka, setCustomParam setCustomParamFunc) (*ClientKafka, error) {
	if kafka == nil {
		return nil, fmt.Errorf("the kafka struct is nil")
	}
	config := sarama.NewConfig()
	//是否开启消息发送成功后通知 successes channel
	config.Producer.Return.Successes = true
	//随机分区器
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 指定消费者的重试策略
	if kafka.ProducerRetryParam != nil {
		config.Producer.Retry.Max = kafka.ProducerRetryParam.Max
		config.Producer.Retry.Backoff = kafka.ProducerRetryParam.Backoff
		config.Producer.Retry.BackoffFunc = kafka.ProducerRetryParam.BackoffFunc
	}
	// 是否需要等待broker的ack
	config.Producer.RequiredAcks = kafka.ProducerAck
	//max.in.flight.requests.per.connection 这个参数没有设置

	// 自定义参数设置
	if setCustomParam != nil {
		setCustomParam(config)
	}
	producerClinet, err := sarama.NewClient(kafka.ServerAddr, config)
	return &ClientKafka{ProducerClient: producerClinet}, err
}

func newKafkaConsumerClient(kafka *Kafka, setCustomParam setCustomParamFunc) (*ClientKafka, error) {
	if kafka == nil {
		return nil, fmt.Errorf("the kafka struct is nil")
	}
	config := sarama.NewConfig()
	// 开启消费报错，将错误投递到错误的channel中
	config.Consumer.Return.Errors = true
	// 消费消息的重试参数
	if kafka.ConsumerOffsetParam != nil {
		config.Consumer.Retry.Backoff = kafka.ConsumerOffsetParam.RetryBackoff
		config.Consumer.Retry.BackoffFunc = kafka.ConsumerOffsetParam.RetryBackoffFunc
		config.Consumer.Offsets.Initial = kafka.ConsumerOffsetParam.ConsumerDefaultOffset
		config.Consumer.Offsets.Retry.Max = kafka.ConsumerOffsetParam.MAX
	}
	//自定义参数
	if setCustomParam != nil {
		setCustomParam(config)
	}
	consumer, err := sarama.NewClient(kafka.ServerAddr, config)
	return &ClientKafka{ConsumerClient: consumer, ConsumerGroupClient: consumer}, err
}

func (client ClientKafka) SendMessage(data *sarama.ProducerMessage) (error, int32, int64) {
	if client.ProducerClient == nil {
		return fmt.Errorf("please call ReceiveMessageByConsumer after initialization NewKafkaProductClient"), 0, 0
	}
	producer, err := sarama.NewSyncProducerFromClient(client.ProducerClient)
	if err != nil {
		return err, 0, 0
	}
	partition, offset, err := producer.SendMessage(data)
	if err != nil {
		return err, 0, 0
	}
	return err, partition, offset
}

func (client ClientKafka) ReceiveMessageByConsumer(topic string, callBack func(data *sarama.ConsumerMessage,
	ch <-chan *sarama.ConsumerError)) error {
	if client.ConsumerClient == nil {
		return fmt.Errorf("please call ReceiveMessageByConsumer after initialization NewKafkaConsumerClient")
	}
	var wg sync.WaitGroup
	c, err := sarama.NewConsumerFromClient(client.ConsumerClient)
	if err != nil {
		return err
	}
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
		// 获取消费error的消息
		var errorCh = make(chan *sarama.ConsumerError)
		go func() {
			ch := <-pc.Errors()
			errorCh <- ch
		}()
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				callBack(msg, errorCh)
			}
		}(pc)
	}
	wg.Wait()
	defer func() {
		_ = c.Close()
	}()
	return nil
}

func (client ClientKafka) ReceiveMessageByConsumerGroup(groupId string, context context.Context,
	topics []string, handler sarama.ConsumerGroupHandler) error {
	if client.ConsumerClient == nil {
		return fmt.Errorf("please call ReceiveMessageByConsumer after initialization NewKafkaConsumerClient")
	}
	cGroup, err := sarama.NewConsumerGroupFromClient(groupId, client.ConsumerGroupClient)
	if err != nil {
		return err
	}
	err = cGroup.Consume(context, topics, handler)
	return err
}

// 下面的就是一个例子
type ConsumerByGroup struct{}

func (ConsumerByGroup) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerByGroup) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (consumer ConsumerByGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var consumerCount int
	for msg := range claim.Messages() {
		time.Sleep(time.Second * 1)
		fmt.Printf("Message topic:%q partition:%d offset:%d  value:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		// 插入mysql
		//TODO 业务逻辑
		/*去做业务逻辑*/
		// 手动提交模式下，也需要先进行标记
		sess.MarkMessage(msg, "")
		consumerCount++
		if consumerCount%3 == 0 {
			// 手动提交，不能频繁调用
			t1 := time.Now().Nanosecond()
			sess.Commit()
			t2 := time.Now().Nanosecond()
			fmt.Println("commit cost:", (t2-t1)/(1000*1000), "ms")
		}
	}
	return nil
}
