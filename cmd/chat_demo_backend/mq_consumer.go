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
		cfg.RabbitMQ.MsgQueueDeclare.Name,
		cfg.RabbitMQ.MsgQueueDeclare.Durable,
		cfg.RabbitMQ.MsgQueueDeclare.DeleteWhenUnused,
		cfg.RabbitMQ.MsgQueueDeclare.Exclusive,
		cfg.RabbitMQ.MsgQueueDeclare.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}
	m.q = q

	err = ch.QueueBind(
		q.Name,
		cfg.RabbitMQ.MsgQueueBind.RoutingKey,
		cfg.RabbitMQ.MsgQueueBind.Exchange,
		cfg.RabbitMQ.MsgQueueBind.NoWait,
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
		m.cfg.RabbitMQ.MsgConsume.Consumer,
		m.cfg.RabbitMQ.MsgConsume.AutoAck,
		m.cfg.RabbitMQ.MsgConsume.Exclusive,
		m.cfg.RabbitMQ.MsgConsume.NoLocal,
		m.cfg.RabbitMQ.MsgConsume.NoWait,
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
