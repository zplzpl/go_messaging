package main

import (
	"github.com/streadway/amqp"
	"go_messaging/pkg/logger"
)

type MQPublisher struct {
	cfg *ConfigParser
	ch  *amqp.Channel
}

func NewMQPublisher(cfg *ConfigParser, ch *amqp.Channel) *MQPublisher {
	return &MQPublisher{cfg: cfg, ch: ch}
}

func (a *MQPublisher) Publish(msgKey string, buf []byte) error {

	err := a.ch.Publish(
		a.cfg.RabbitMQ.Publish.Exchange,
		msgKey,
		a.cfg.RabbitMQ.Publish.Mandatory,
		a.cfg.RabbitMQ.Publish.Immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        buf,
		},
	)

	if err != nil {
		logger.GetLogger().Error("Failed to publish a message", logger.Error(err))
		return err
	}

	return nil
}
