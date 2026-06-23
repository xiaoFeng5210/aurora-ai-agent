package rabbitmq_module

import (
	"aurora-agent/utils"

	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"go.uber.org/zap"
)

var (
	Conn       *amqp.Connection
	ConnectErr error
	once       sync.Once
)

func Connect() (*amqp.Connection, error) {
	conf := utils.InitViper("conf", "rabbitmq", "yaml")
	url := conf.GetString("rabbitmq.url")
	conn, ConnectErr := amqp.Dial(url)
	if ConnectErr != nil {
		utils.Logger.Error("Failed to connect to rabbitmq.", zap.Error(ConnectErr))
		return nil, ConnectErr
	}
	utils.Logger.Info("Connected to rabbitmq successfully.")

	return conn, ConnectErr
}
