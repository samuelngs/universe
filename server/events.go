package server

import "encoding/json"

// Event topics
const (
	EventServerStart   string = "server-start"
	EventServerStop           = "server-stop"
	EventServerStarted        = "server-started"
	EventServerStopped        = "server-stopped"
	EventReceiveSignal        = "receive-signal"
)

// Event interface for secure shell server
type Event interface {
	Topic() string
	Message() string
	String() string
}

// secure shell event
type event struct {
	topic, message string
}

// Topic returns event topic
func (v *event) Topic() string {
	return v.topic
}

// Message returns trace event message
func (v *event) Message() string {
	return v.message
}

// String returns event object in string format
func (v *event) String() string {
	o := map[string]interface{}{
		"topic": v.topic,
	}
	if v.message != "" {
		o["message"] = v.message
	}
	b, _ := json.Marshal(o)
	return string(b[:])
}

// ParseEvent message and return trace object
func ParseEvent(s string) Event {
	o := new(event)
	json.Unmarshal([]byte(s), o)
	return o
}
