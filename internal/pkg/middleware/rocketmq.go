package middleware

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/spf13/viper"
)

type Rocketmq struct {
	provider    *rocketmq.Producer
	consumer    *rocketmq.PushConsumer
	nameSrvAddr string
	syncData    SyncData
}

type SyncData struct {
	topic     string
	groupName string
	retryTime int
	tag       string
}

// NewRocketmqProvider create the producer of rocketmq.
func NewRocketmqProvider(conf *viper.Viper) (rocketmq.Producer, error) {
	nameSrvAddr := conf.GetString("data.rocketmq.namesrvAddr")
	topic := conf.GetString("data.rocketmq.syncdata.topic")
	groupName := conf.GetString("data.rocketmq.syncdata.groupName")
	retryTime := conf.GetInt("data.rockertmq.syncdata.retryTime")
	srvAddr, err := primitive.NewNamesrvAddr(nameSrvAddr)

	p, err := rocketmq.NewProducer(
		producer.WithGroupName(groupName),
		producer.WithNameServer(srvAddr),
		producer.WithCreateTopicKey(topic),
		producer.WithRetry(retryTime),
	)

	err = p.Start()
	return p, err
}

func NewRocketmqConsumer(conf *viper.Viper) rocketmq.PushConsumer {
	nameSrvaddr := conf.GetString("data.rocketmq.namesrvAddr")
	groupName := conf.GetString("data.rocketmq.syncdata.groupName")
	retryTime := conf.GetInt("data.rockertmq.syncdata.retryTime")
	srvAddr, err := primitive.NewNamesrvAddr(nameSrvaddr)

	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(srvAddr),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName(groupName),
		consumer.WithRetry(retryTime),
	)
	if err != nil {
		panic(err)
	}
	return c
}

// SendMessage send message with producer and return the result using the type of boolean.
func SendMessage(p rocketmq.Producer, topic, content string) bool {
	msg := &primitive.Message{
		Topic: topic,
		Body:  []byte(content),
	}
	result, err := p.SendSync(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	status := result.Status
	return status == primitive.SendOK
}

func SendMessageWithTag(p rocketmq.Producer, topic, content, tag string) bool {
	msg := &primitive.Message{
		Topic: topic,
		Body:  []byte(content),
	}
	msg.WithTag(tag)
	result, err := p.SendSync(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	status := result.Status
	return status == primitive.SendOK
}

func SendMessageBatch(p rocketmq.Producer, topic string, content []string) bool {
	var msgs []*primitive.Message
	for _, ctx := range content {
		msgs = append(msgs, &primitive.Message{
			Topic: topic,
			Body:  []byte(ctx),
		})
	}
	result, err := p.SendSync(context.Background(), msgs...)
	if err != nil {
		panic(err)
	}
	status := result.Status
	return status == primitive.SendOK
}

// SubMessage subscribe the message with the rocket consumer.
func SubMessage(c rocketmq.PushConsumer, topic string) []string {
	var body []string
	err := c.Subscribe(topic, consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "*",
	}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		orderlyCtx, _ := primitive.GetOrderlyCtx(ctx)
		fmt.Printf("orderly context: %v\n", orderlyCtx)
		for i := range ext {
			body = append(body, string(ext[i].Body))
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		panic(err)
	}
	return body
}
