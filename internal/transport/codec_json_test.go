package transport

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	dino "github.com/angyxys/micro-dino"
)

func TestJSONCodec_Name(t *testing.T) {
	c := NewJSONCodec()
	if got, want := c.Name(), "json"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestJSONCodec_MarshalNil(t *testing.T) {
	c := NewJSONCodec()
	if _, err := c.Marshal(nil); err == nil {
		t.Fatal("expected error marshaling nil message")
	}
}

func TestJSONCodec_UnmarshalNilTarget(t *testing.T) {
	c := NewJSONCodec()
	err := c.Unmarshal([]byte(`{}`), nil)
	if err == nil {
		t.Fatal("expected error unmarshaling into nil target")
	}
}

func TestJSONCodec_UnmarshalEmptyPayload(t *testing.T) {
	c := NewJSONCodec()
	err := c.Unmarshal(nil, &Message{})
	if err == nil {
		t.Fatal("expected error on empty payload")
	}
}

func TestJSONCodec_UnmarshalInvalidJSON(t *testing.T) {
	c := NewJSONCodec()
	err := c.Unmarshal([]byte(`{not-json`), &Message{})
	if err == nil {
		t.Fatal("expected error on invalid JSON")
	}
}

func TestJSONCodec_RoundTrip_Request(t *testing.T) {
	c := NewJSONCodec()
	orig := NewRequest("users.get", dino.Params{
		"id":    "u-123",
		"limit": 50,
		"tags":  []any{"admin", "active"},
	})
	orig.Meta = dino.Meta{"tenant": "acme"}
	orig.TraceID = "trace-abc"
	orig.NodeID = "node-1"
	orig.Service = "users"
	orig.Version = "v1"

	data, err := c.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got Message
	if err := c.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Kind != KindRequest {
		t.Errorf("Kind = %q, want %q", got.Kind, KindRequest)
	}
	if got.ID != orig.ID {
		t.Errorf("ID = %q, want %q", got.ID, orig.ID)
	}
	if got.Action != "users.get" {
		t.Errorf("Action = %q", got.Action)
	}
	if got.Params.Get("id").String() != "u-123" {
		t.Errorf("Params.id = %v", got.Params.Get("id").Any())
	}
	if n, _ := got.Params.Get("limit").Int64(); n != 50 {
		t.Errorf("Params.limit = %v", got.Params.Get("limit").Any())
	}
	if got.Meta.Get("tenant").String() != "acme" {
		t.Errorf("Meta.tenant = %v", got.Meta.Get("tenant").Any())
	}
	if got.TraceID != "trace-abc" {
		t.Errorf("TraceID = %q", got.TraceID)
	}
	if got.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

func TestJSONCodec_RoundTrip_Response_Success(t *testing.T) {
	c := NewJSONCodec()
	req := NewRequest("users.get", nil)
	orig := NewResponse(req, map[string]any{"id": "u-1", "name": "Ana"}, nil)

	data, err := c.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got Message
	if err := c.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Kind != KindResponse {
		t.Errorf("Kind = %q", got.Kind)
	}
	if got.ReplyTo != req.ID {
		t.Errorf("ReplyTo = %q, want %q", got.ReplyTo, req.ID)
	}
	if got.Error != nil {
		t.Errorf("unexpected Error: %v", got.Error)
	}

	result, ok := got.Result.(map[string]any)
	if !ok {
		t.Fatalf("Result type = %T, want map[string]any", got.Result)
	}
	if result["name"] != "Ana" {
		t.Errorf("Result.name = %v", result["name"])
	}
}

func TestJSONCodec_RoundTrip_Response_Error(t *testing.T) {
	c := NewJSONCodec()
	req := NewRequest("users.get", nil)
	orig := NewResponse(req, nil, &WireError{
		Code:    "NOT_FOUND",
		Message: "user not found",
		Details: map[string]any{"id": "u-1"},
	})

	data, err := c.Marshal(orig)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got Message
	if err := c.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Error == nil {
		t.Fatal("expected Error, got nil")
	}
	if got.Error.Code != "NOT_FOUND" {
		t.Errorf("Code = %q", got.Error.Code)
	}
	if got.Error.Message != "user not found" {
		t.Errorf("Message = %q", got.Error.Message)
	}
	if got.Error.Error() != "user not found" {
		t.Errorf("Error() = %q", got.Error.Error())
	}
}

func TestJSONCodec_PreservesLargeIntegers(t *testing.T) {
	c := NewJSONCodec()
	orig := NewRequest("x", dino.Params{"big": int64(9007199254740993)})

	data, _ := c.Marshal(orig)
	var got Message
	if err := c.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	raw := got.Params.Get("big").Any()
	if _, ok := raw.(json.Number); !ok {
		t.Fatalf("expected json.Number, got %T", raw)
	}
	n, err := got.Params.Get("big").Int64()
	if err != nil {
		t.Fatalf("Int64: %v", err)
	}
	if n != 9007199254740993 {
		t.Errorf("precision lost: got %d", n)
	}
}

func TestJSONCodec_RoundTrip_Event(t *testing.T) {
	c := NewJSONCodec()
	orig := &Message{
		Kind:      KindEvent,
		ID:        NewID(),
		Action:    "user.created",
		Params:    dino.Params{"id": "u-1"},
		Timestamp: time.Now().UTC(),
	}

	data, _ := c.Marshal(orig)
	var got Message
	if err := c.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Kind != KindEvent {
		t.Errorf("Kind = %q", got.Kind)
	}
	if got.Action != "user.created" {
		t.Errorf("Action = %q", got.Action)
	}
}

func TestMessage_Validate(t *testing.T) {
	cases := []struct {
		name    string
		msg     *Message
		wantErr bool
	}{
		{"nil", nil, true},
		{"empty kind", &Message{ID: "1"}, true},
		{"empty id", &Message{Kind: KindRequest, Action: "x"}, true},
		{"request no action", &Message{Kind: KindRequest, ID: "1"}, true},
		{"valid request", NewRequest("users.get", nil), false},
		{"valid event", &Message{Kind: KindEvent, ID: "1"}, false},
		{"unknown kind", &Message{Kind: "bogus", ID: "1"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate() err = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestNewID_UniqueAndFormatted(t *testing.T) {
	seen := make(map[string]struct{}, 1000)
	for i := 0; i < 1000; i++ {
		id := NewID()
		if len(id) != 24 {
			t.Fatalf("id length = %d, want 24", len(id))
		}
		if strings.ContainsAny(id, "ghijklmnopqrstuvwxyz") {
			t.Fatalf("id not hex: %q", id)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id: %q", id)
		}
		seen[id] = struct{}{}
	}
}

func TestCodecRegistry(t *testing.T) {
	if _, ok := GetCodec(CodecJSONName); !ok {
		t.Fatal("json codec not registered by init()")
	}

	if c := MustCodec(CodecJSONName); c.Name() != CodecJSONName {
		t.Fatalf("MustCodec returned wrong codec: %q", c.Name())
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unknown codec")
		}
	}()
	MustCodec("does-not-exist")
}
