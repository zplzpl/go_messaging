package main

import (
	"github.com/streadway/amqp"
)

type ConsumeHandler interface {
	ConsumeMsg(msg []byte) error
}

type MQConsumer struct {
	cfg *ConfigParser
	h   ConsumeHandler
	ch  *amqp.Channel
	q   amqp.Queue
}

func NewMQConsumer(cfg *ConfigParser, ch *amqp.Channel, h ConsumeHandler) (*MQConsumer, error) {

	m := &MQConsumer{cfg: cfg, ch: ch, h: h}

	q, err := ch.QueueDeclare(
		cfg.RabbitMQ.QueueDeclare.Name,
		cfg.RabbitMQ.QueueDeclare.Durable,
		cfg.RabbitMQ.QueueDeclare.DeleteWhenUnused,
		cfg.RabbitMQ.QueueDeclare.Exclusive,
		cfg.RabbitMQ.QueueDeclare.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}
	m.q = q

	err = ch.QueueBind(
		q.Name,
		cfg.RabbitMQ.QueueBind.RoutingKey,
		cfg.RabbitMQ.QueueBind.Exchange,
		cfg.RabbitMQ.QueueBind.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return m, nil

}

func (m *MQConsumer) Run() error {

	msgs, err := m.ch.Consume(
		m.q.Name,
		m.cfg.RabbitMQ.Consume.Consumer,
		m.cfg.RabbitMQ.Consume.AutoAck,
		m.cfg.RabbitMQ.Consume.Exclusive,
		m.cfg.RabbitMQ.Consume.NoLocal,
		m.cfg.RabbitMQ.Consume.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case d := <-msgs:
				_ = m.h.ConsumeMsg(d.Body)
			}
		}
	}()

	return nil

}
