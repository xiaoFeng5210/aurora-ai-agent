package rabbitmq_module

import (
	"aurora-agent/utils"

	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"go.uber.org/zap"
)

var (
	Conn *amqp.Connection
	ConnectErr error
	once sync.Once
)

func Connect() (*amqp.Connection, error) {
	once.Do(func() {
		conf := utils.InitViper("conf", "rabbitmq", "yaml")
		url := conf.GetString("rabbitmq.url")
		Conn, ConnectErr = amqp.Dial(url)
		if ConnectErr != nil {
			utils.Logger.Error("Failed to connect to rabbitmq.", zap.Error(ConnectErr))
			return
		}
		utils.Logger.Info("Connected to rabbitmq successfully.")
	})

	return Conn, ConnectErr
}


