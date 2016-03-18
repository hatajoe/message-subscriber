package pubsub

import "github.com/garyburd/redigo/redis"

type PubSub struct {
	pubsubConn *redis.PubSubConn
	channel    string
	buffer     chan []byte
}

func NewPubSub(conn redis.Conn, channel string, buffer chan []byte) (*PubSub, error) {
	psc := redis.PubSubConn{conn}
	if err := psc.Subscribe(channel); err != nil {
		return nil, err
	}
	return &PubSub{
		pubsubConn: &psc,
		channel:    channel,
		buffer:     buffer,
	}, nil
}

func (m *PubSub) Subscribe() error {
	switch v := m.pubsubConn.Receive().(type) {
	case redis.Message:
		m.buffer <- v.Data
	case error:
		return v
	default:
	}
	return nil
}

func (m *PubSub) UnSubscribe() error {
	if err := m.pubsubConn.Unsubscribe(m.channel); err != nil {
		return err
	}
	if err := m.pubsubConn.Close(); err != nil {
		return err
	}
	return nil
}

func (m *PubSub) Abort() error {
	return nil
}

func (m *PubSub) End() error {
	return nil
}
