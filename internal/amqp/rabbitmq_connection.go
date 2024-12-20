package amqp

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func CreateConnectionRabbitmq() *amqp.Channel {

	host := viper.GetString("rabbitmq.host")
	user := viper.GetString("rabbitmq.user")
	password := viper.GetString("rabbitmq.password")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:5672", user, password, host))

	failOnError(err, "failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return ch
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
