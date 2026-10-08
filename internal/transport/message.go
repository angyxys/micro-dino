package transport

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	dino "github.com/angyxys/micro-dino"
)

type Kind string

const (
	KindRequest    Kind = "req"
	KindResponse   Kind = "res"
	KindEvent      Kind = "evt"
	KindHeartbeat  Kind = "hbt"
	KindRegister   Kind = "reg"
	KindDeregister Kind = "dreg"
)

type Message struct {
	Kind    Kind   `json:"kind"`
	ID      string `json:"id"`
	ReplyTo string `json:"reply_to,omitempty"`

	Action string      `json:"action,omitempty"` // req
	Params dino.Params `json:"params,omitempty"` // req
	Result any         `json:"result,omitempty"` // res (éxito)
	Error  *WireError  `json:"error,omitempty"`  // res (error)

	Meta    dino.Meta `json:"meta,omitempty"`
	TraceID string    `json:"trace_id,omitempty"`

	NodeID    string    `json:"node_id,omitempty"`
	Service   string    `json:"service,omitempty"`
	Version   string    `json:"version,omitempty"`
	Timestamp time.Time `json:"ts"`
}

type WireError struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func (e *WireError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func NewWireError(code string, err error) *WireError {
	if err == nil {
		return nil
	}
	return &WireError{Code: code, Message: err.Error()}
}

func NewID() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {

		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b[:])
}

func (m *Message) Validate() error {
	if m == nil {
		return errors.New("micro-dino: nil message")
	}
	if m.Kind == "" {
		return errors.New("micro-dino: message kind is empty")
	}
	if m.ID == "" {
		return errors.New("micro-dino: message id is empty")
	}
	switch m.Kind {
	case KindRequest, KindResponse:

	case KindEvent, KindHeartbeat, KindRegister, KindDeregister:

	default:
		return errors.New("micro-dino: unknown message kind: " + string(m.Kind))
	}
	if m.Kind == KindRequest && m.Action == "" {
		return errors.New("micro-dino: request without action")
	}
	return nil
}

func NewRequest(action string, params dino.Params) *Message {
	return &Message{
		Kind:      KindRequest,
		ID:        NewID(),
		Action:    action,
		Params:    params,
		Timestamp: time.Now().UTC(),
	}
}

func NewResponse(req *Message, result any, err *WireError) *Message {
	return &Message{
		Kind:      KindResponse,
		ID:        NewID(),
		ReplyTo:   req.ID,
		Result:    result,
		Error:     err,
		Timestamp: time.Now().UTC(),
	}
}
