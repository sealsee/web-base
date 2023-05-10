package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Shopify/sarama"
	"github.com/sealsee/web-base/public/ds"
	"go.uber.org/zap"
)

var prd sarama.SyncProducer
var com sarama.ConsumerGroup

func Init() {
	prd = ds.GetKafkaPrd()
	com = ds.GetKafkaCom()
}

type Handler func(string) error

func Send(topic string, a any) {
	if topic == "" || a == nil {
		return
	}
	kind := reflect.TypeOf(a).Kind()
	if kind != reflect.Map && kind != reflect.Struct {
		return
	}

	bytes, _ := json.Marshal(a)
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(bytes)
	pid, offset, err := prd.SendMessage(msg)
	if err != nil {
		zap.L().Error("kafka send:", zap.Error(err))
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)

}

type consumerGroupHandler struct {
	Handler Handler
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		sess.MarkMessage(msg, "")
		fmt.Println(string(msg.Value))
	}
	return nil
}

func Receive(topic string, handler Handler) {
	if topic == "" || handler == nil {
		return
	}

	kind := reflect.TypeOf(handler).Kind()
	if kind != reflect.Func {
		panic("handler is not func")
	}

	go func() {
		for {
			com.Consume(context.Background(), []string{topic}, &consumerGroupHandler{Handler: handler})
		}
	}()
}
