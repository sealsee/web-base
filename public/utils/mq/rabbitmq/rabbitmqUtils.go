package rabbitmq

import (
	"encoding/json"
	"reflect"

	"github.com/sealsee/web-base/public/ds"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var chn *amqp.Channel

func Init() {
	chn = ds.GetRabbitMQChn()
}

type Handler func(string) error

func Declare(exchange, key string) {
	if exchange == "" || key == "" {
		return
	}
	chn.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	chn.QueueDeclare(key, true, false, false, false, nil)
	chn.QueueBind(key, key, exchange, false, nil)
}

func Send(exchange, key string, a any) {
	if exchange == "" || key == "" || a == nil {
		return
	}

	kind := reflect.TypeOf(a).Kind()
	if kind != reflect.Map && kind != reflect.Struct {
		return
	}

	bytes, _ := json.Marshal(a)
	err := chn.Publish(exchange, key, false, false, amqp.Publishing{ContentType: "text/plain", Body: bytes})
	if err != nil {
		zap.L().Error("rabbitmq send:", zap.Error(err))
	}
}

func Receive(queue string, handler Handler) {
	if queue == "" || handler == nil {
		return
	}

	kind := reflect.TypeOf(handler).Kind()
	if kind != reflect.Func {
		panic("handler is not func")
	}

	msgs, err := chn.Consume(queue, "", false, false, false, false, nil)
	// chn.Qos(1, 1, true)
	if err != nil {
		zap.L().Error("rabbitmq receive:", zap.Error(err))
	}
	go func() {
		for msg := range msgs {
			err := handler(string(msg.Body))
			if err == nil {
				msg.Ack(false)
			} else {
				zap.L().Error("rabbitmq receive handle:", zap.Error(err))
			}
		}
	}()
}
