package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func produce(msg string, ch *amqp.Channel, exg, key string) {
}
