package server

import "encoding/json"

// Log topics
const (
	TraceAuthentication                    string = "authentication"
	TracePasswordAuthentication                   = "password-authentication"
	TraceRSAAuthentication                        = "rsa-authentication"
	TraceKeyboardInteractiveAuthentication        = "keyboardinteractive-authentication"
	TraceHandshake                                = "handshake"
	TraceChannel                                  = "channel"
	TraceConnect                                  = "connect"
	TraceDisconnect                               = "disconnect"
)

// Log interface
type Log interface {
	Topic() string
	Error() error
	Message() string
	String() string
}

// secure shell trace log
type trace struct {
	topic, message string
	err            error
}

// Topic returns trace log topic
func (v *trace) Topic() string {
	return v.topic
}

// Error returns trace log error
func (v *trace) Error() error {
	return v.err
}

// Message returns trace log message
func (v *trace) Message() string {
	return v.message
}

// String returns trace log object in string format
func (v *trace) String() string {
	o := map[string]interface{}{
		"topic": v.topic,
	}
	if v.message != "" {
		o["message"] = v.message
	}
	if v.err != nil {
		o["error"] = v.err
	}
	b, _ := json.Marshal(o)
	return string(b[:])
}

// ParseLog message and return trace object
func ParseLog(s string) Log {
	o := new(trace)
	err := json.Unmarshal([]byte(s), o)
	if err != nil {
		o.err = err
	}
	return o
}
