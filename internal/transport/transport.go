package transport

import (
	"context"
)

type Transport interface {
	Name() string

	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Connected() bool

	Request(ctx context.Context, target Node, msg *Message) (*Message, error)

	Send(ctx context.Context, target Node, msg *Message) error

	Broadcast(ctx context.Context, topic string, msg *Message) error

	Subscribe(handler Handler) error
}

type Handler func(ctx context.Context, msg *Message) (*Message, error)

type Node struct {
	ID      string
	Address string
	Service string
	Version string
	Meta    map[string]string
}
