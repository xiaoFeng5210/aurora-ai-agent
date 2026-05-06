package rabbitmq

import (
	"aurora-agent/utils"

	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"go.uber.org/zap"
)

var (
	conn *amqp.Connection
	connectErr error
	once sync.Once
)

func Connect() (*amqp.Connection, error) {
	once.Do(func() {
		conf := utils.InitViper("conf", "rabbitmq", "yaml")

		url := conf.GetString("rabbitmq.url")
		conn, connectErr = amqp.Dial(url)
		if connectErr != nil {
			utils.Logger.Error("Failed to connect to rabbitmq.", zap.Error(connectErr))
			return
		}
		utils.Logger.Info("Connected to rabbitmq successfully.")
	})

	return conn, connectErr
}
