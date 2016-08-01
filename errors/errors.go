package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// DefaultRef string
const DefaultRef string = "-"

// Error provide a way to return detailed information
// for an error. The error is normally JSON encoded.
type Error struct {
	Ref    string `json:"resource"`
	Code   int    `json:"code"`
	Reason string `json:"reason"`
	Detail string `json:"detail,omitempty"`
	Status string `json:"status"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b[:])
}

// Info add additional resource info to error
func (e *Error) Info(s interface{}) *Error {
	switch v := s.(type) {
	case string:
		e.Detail = v
	case error:
		e.Detail = v.Error()
	case int:
		e.Detail = strconv.Itoa(v)
	case int8:
		e.Detail = strconv.Itoa(int(v))
	case int16:
		e.Detail = strconv.Itoa(int(v))
	case int32:
		e.Detail = strconv.Itoa(int(v))
	case int64:
		e.Detail = strconv.FormatInt(v, 10)
	case uint:
		e.Detail = strconv.FormatUint(uint64(v), 10)
	case uint8:
		e.Detail = strconv.FormatUint(uint64(v), 10)
	case uint16:
		e.Detail = strconv.FormatUint(uint64(v), 10)
	case uint32:
		e.Detail = strconv.FormatUint(uint64(v), 10)
	case uint64:
		e.Detail = strconv.FormatUint(uint64(v), 10)
	case float32:
		e.Detail = strconv.FormatFloat(float64(v), 'f', 6, 64)
	case float64:
		e.Detail = strconv.FormatFloat(v, 'f', 6, 64)
	case bool:
		e.Detail = strconv.FormatBool(v)
	default:
		e.Detail = fmt.Sprintf("%v", s)
	}
	return e
}

// New creates a new error
func New(ref, reason string, code int) *Error {
	return &Error{
		Ref:    ref,
		Code:   code,
		Reason: reason,
		Status: http.StatusText(int(code)),
	}
}

// Parse to read error message
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Reason = err
	}
	return e
}

// Read error
func Read(e error) *Error {
	switch o := e.(type) {
	case *Error:
		return o
	default:
		return Forbidden(DefaultRef, e.Error())
	}
}

// BadRequest error
func BadRequest(ref, reason string) *Error {
	return &Error{
		Ref:    ref,
		Code:   http.StatusBadRequest,
		Reason: reason,
		Status: http.StatusText(http.StatusBadRequest),
	}
}

// Unauthorized error
func Unauthorized(ref, reason string) *Error {
	return &Error{
		Ref:    ref,
		Code:   http.StatusUnauthorized,
		Reason: reason,
		Status: http.StatusText(http.StatusUnauthorized),
	}
}

// Forbidden error
func Forbidden(ref, reason string) *Error {
	return &Error{
		Ref:    ref,
		Code:   http.StatusForbidden,
		Reason: reason,
		Status: http.StatusText(http.StatusForbidden),
	}
}

// NotFound error
func NotFound(ref, reason string) *Error {
	return &Error{
		Ref:    ref,
		Code:   http.StatusNotFound,
		Reason: reason,
		Status: http.StatusText(http.StatusNotFound),
	}
}

// InternalServer error
func InternalServer(ref, reason string) *Error {
	return &Error{
		Ref:    ref,
		Code:   http.StatusInternalServerError,
		Reason: reason,
		Status: http.StatusText(http.StatusInternalServerError),
	}
}
