package discovery

import (
	"context"

	"github.com/angyxys/micro-dino/internal/transport"
)

type Discovery interface {
	Register(ctx context.Context, self transport.Node) error

	Deregister(ctx context.Context, self transport.Node) error

	Lookup(ctx context.Context, action string, filter Filter) ([]transport.Node, error)

	Watch(ctx context.Context, action string) (<-chan Event, error)

	Close() error
}

type Filter struct {
	Service string
	Version string
	Meta    map[string]string
}

type Event struct {
	Type NodeEventType
	Node transport.Node
}

type NodeEventType string

const (
	NodeAdded   NodeEventType = "added"
	NodeRemoved NodeEventType = "removed"
	NodeUpdated NodeEventType = "updated"
)
