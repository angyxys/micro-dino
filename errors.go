package dino

import "errors"

var (
	ErrServiceNotFound = errors.New("micro-dino: service not found")
	ErrActionNotFound  = errors.New("micro-dino: action not found")
	ErrNoNodes         = errors.New("micro-dino: no nodes available")
	ErrTimeout         = errors.New("micro-dino: timeout")
	ErrTransportClosed = errors.New("micro-dino: transport closed")
	ErrAlreadyStarted  = errors.New("micro-dino: app already started")
)
