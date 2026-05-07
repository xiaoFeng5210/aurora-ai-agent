package rabbitmq_module

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func produce(msg string, ch *amqp.Channel, exg, key string) {
}

func DeclareExchange(ch *amqp.Channel, exg string) error {
	err := ch.ExchangeDeclare(
		exg,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("declare exchange failed: %w", err)
	}
	return nil
}

