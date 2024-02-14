package messaging

import (
	"errors"
)

var errGatewayNotFound = errors.New("not found gateway for consumer")

type ConsumerGatewayRegister struct {
	TopicID  ConsumerTopicID
	Consumer ConsumerReceiveFunc
}

type consumerGateway struct {
	g map[ConsumerTopicID]ConsumerReceiveFunc
}

type ConsumerGateway interface {
	Registers(gws ...ConsumerGatewayRegister)
	Get(key ConsumerTopicID) (ConsumerReceiveFunc, error)
}

func NewConsumerGateway() ConsumerGateway {
	return &consumerGateway{
		g: make(map[ConsumerTopicID]ConsumerReceiveFunc),
	}
}

func (g *consumerGateway) Registers(gws ...ConsumerGatewayRegister) {
	for _, gw := range gws {
		g.g[gw.TopicID] = gw.Consumer
	}
}

func (g *consumerGateway) Get(key ConsumerTopicID) (ConsumerReceiveFunc, error) {
	gateway, ok := g.g[key]
	if !ok {
		return nil, errGatewayNotFound
	}

	return gateway, nil
}
