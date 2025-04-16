package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/project-template/common/config"
	enp "github.com/project-template/common/encapsulate"
	"github.com/project-template/common/tools"
	"github.com/project-template/errorcode"
	"log"
	"os"
	"time"
)

const (
	groupDefault = "kafka_group_default" // kafka默认分组

	KafkaTopicLog           = "trip_portal_log"        //kafka topic 日志
	KafkaTopicTripAccessLog = "trip_portal_access_log" //kafka topic 网约车接口错误日志
	KafkaTopicTripSmsSign   = "trip_portal_sms_sign"   //kafka topic 报送阿里云做短信签名审批
	//KafkaTopicTemp   = "trip_topic"   //kafka topic 报送阿里云做短信签名审批
)

type KafkaMessage struct {
	SubTopic string
	Val      []byte
	Time     time.Time
}

func OpenKafka(config *config.Config) {
	// 初始化生产者
	initProducer(config)
	// 初始化消费者
	initConsumer(groupDefault, config)
}

func initProducer(config *config.Config) {
	var producerConfig = &kafka.ConfigMap{
		"api.version.request": "true",
		"message.max.bytes":   1000000,
		"linger.ms":           10,
		"retries":             30,
		"retry.backoff.ms":    1000,
		"acks":                "1",
	}
	err := producerConfig.SetKey("bootstrap.servers", config.KafkaConfig.Servers)
	if err != nil {
		fmt.Printf("初始化kafka set bootstrap servers 错误:%v\n", err)
		os.Exit(-1)
	}
	err = producerConfig.SetKey("security.protocol", "plaintext")
	if err != nil {
		fmt.Printf("初始化kafka set security protocol 错误:%v\n", err)
		os.Exit(-1)
	}
	config.KafkaProducer, err = kafka.NewProducer(producerConfig)
	if err != nil {
		fmt.Printf("初始化kafka Producer错误:%v\n", err)
		os.Exit(-1)
	}
}

func initConsumer(groupId string, config *config.Config) {
	var err error
	var consumerConfig = &kafka.ConfigMap{
		"api.version.request":       "true",
		"auto.offset.reset":         "latest",
		"heartbeat.interval.ms":     3000,
		"session.timeout.ms":        30000,
		"max.poll.interval.ms":      120000,
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000,
	}
	err = consumerConfig.SetKey("bootstrap.servers", config.KafkaConfig.Servers)
	if err != nil {
		log.Printf("初始化kafka set bootstrap servers 错误:%v\n", err)
		os.Exit(-1)
	}
	err = consumerConfig.SetKey("security.protocol", "plaintext")
	if err != nil {
		log.Printf("初始化kafka set security protocol错误:%v\n", err)
		os.Exit(-1)
	}
	err = consumerConfig.SetKey("group.id", groupId)
	if err != nil {
		log.Printf("初始化kafka set group id 错误:%v\n", err)
		os.Exit(-1)
	}
	config.KafkaConsumer, err = kafka.NewConsumer(consumerConfig)
	if err != nil {
		log.Printf("初始化kafka Consumer错误:%v\n", err)
		os.Exit(-1)
	}
}

func CloseKafka() {
	config.Info().KafkaProducer.Flush(1000)
	config.Info().KafkaProducer.Close()
	err := config.Info().KafkaConsumer.Close()
	if err != nil {
		fmt.Println("关闭消费者错误", err)
	}
}

func KafkaProduce(topic string, message *KafkaMessage) *enp.Response {
	bytes, err := json.Marshal(message)
	if err != nil {
		return enp.Put(errorcode.JsonMarshal, enp.AddIn(message), enp.AddError(err))
	}
	tools.SecureGo(func(args ...interface{}) {
		// 监听消息发送结果
		for e := range config.Info().KafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					// 发送失败
					enp.Put(errorcode.KafkaProduceTopicPartitionError, enp.AddError(ev.TopicPartition.Error))
				} else {
					// 发送成功.......
				}
			}
		}
	})
	err = config.Info().KafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          bytes,
	}, nil)
	if err != nil {
		return enp.Put(errorcode.KafkaProduce, enp.AddError(err))
	}

	return enp.Put(errorcode.Success)
}

// KafkaConsumer group使用场景说明：
// 本系统内完全不考虑多次多个group消费一个topic的情况，每个topic都使用默认的group消费
func KafkaConsumer(topic string, f func(message []byte) *enp.Response) *enp.Response {
	var err error
	err = config.Info().KafkaConsumer.Subscribe(topic, nil)
	if err != nil {
		return enp.Put(errorcode.KafkaConsumerSubscribe, enp.AddError(err))
	}
	tools.SecureGo(func(args ...interface{}) {
		for {
			var message *kafka.Message
			message, err = config.Info().KafkaConsumer.ReadMessage(-1)
			if err != nil || message == nil {
				enp.Put(errorcode.KafkaConsumerReadMessage, enp.AddError(err))
				return
			}
			f(message.Value)
			_, err = config.Info().KafkaConsumer.Commit()
			if err != nil {
				enp.Put(errorcode.KafkaCommit, enp.AddError(err))
				return
			}
		}
	})

	return enp.Put(errorcode.Success)
}
