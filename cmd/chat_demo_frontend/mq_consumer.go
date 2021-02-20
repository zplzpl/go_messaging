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
		cfg.RabbitMQ.MsgRespQueueDeclare.Name,
		cfg.RabbitMQ.MsgRespQueueDeclare.Durable,
		cfg.RabbitMQ.MsgRespQueueDeclare.DeleteWhenUnused,
		cfg.RabbitMQ.MsgRespQueueDeclare.Exclusive,
		cfg.RabbitMQ.MsgRespQueueDeclare.NoWait,
		nil,
	)
	if err != nil {
		return nil, err
	}
	m.q = q

	err = ch.QueueBind(
		q.Name,
		cfg.RabbitMQ.MsgRespQueueBind.RoutingKey,
		cfg.RabbitMQ.MsgRespQueueBind.Exchange,
		cfg.RabbitMQ.MsgRespQueueBind.NoWait,
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
		m.cfg.RabbitMQ.MsgRespConsume.Consumer,
		m.cfg.RabbitMQ.MsgRespConsume.AutoAck,
		m.cfg.RabbitMQ.MsgRespConsume.Exclusive,
		m.cfg.RabbitMQ.MsgRespConsume.NoLocal,
		m.cfg.RabbitMQ.MsgRespConsume.NoWait,
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
