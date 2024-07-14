package errors

import (
	"encoding/json"
	"fmt"
)

type ErrorCode int

const (
	GenericErrorCode   ErrorCode = 1
	NotImplementedCode ErrorCode = 2
)

type ErrorWithCode interface {
	Code() ErrorCode
}

type errorImpl struct {
	Value   ErrorCode `json:"Code"`
	Message string
	Cause   error `json:",omitempty"`
}

func New(message string) error {
	return errorImpl{
		Value:   GenericErrorCode,
		Message: message,
	}
}

func NewCode(code ErrorCode) error {
	e := errorImpl{
		Value:   code,
		Message: determineCommonErrorMessage(code),
	}

	return e
}

func NotImplemented() error {
	return NewCode(NotImplementedCode)
}

func NewCodeWithDetails(code ErrorCode, details string) error {
	e := errorImpl{
		Value:   code,
		Message: details,
	}

	return e
}

func Newf(format string, args ...interface{}) error {
	return errorImpl{
		Value:   GenericErrorCode,
		Message: fmt.Sprintf(format, args...),
	}
}

func Wrap(cause error, message string) error {
	return errorImpl{
		Value:   GenericErrorCode,
		Message: message,
		Cause:   cause,
	}
}

func WrapCode(cause error, code ErrorCode) error {
	e := errorImpl{
		Value:   code,
		Message: determineCommonErrorMessage(code),
		Cause:   cause,
	}

	return e
}

func Wrapf(cause error, format string, args ...interface{}) error {
	return errorImpl{
		Value:   GenericErrorCode,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	ie, ok := err.(errorImpl)
	if !ok {
		return nil
	}

	return ie.Cause
}

func IsErrorWithCode(err error, code ErrorCode) bool {
	if err == nil {
		return false
	}

	if impl, ok := err.(errorImpl); ok {
		return impl.Code() == code
	}

	return false
}

func (e errorImpl) Error() string {
	var out string

	out += fmt.Sprintf("(%d) ", e.Value)
	out += e.Message

	if e.Cause != nil {
		out += fmt.Sprintf(" (cause: %v)", e.Cause.Error())
	}

	return out
}

func (e errorImpl) Code() ErrorCode {
	return e.Value
}

func (e errorImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    ErrorCode
		Message string          `json:",omitempty"`
		Cause   json.RawMessage `json:",omitempty"`
	}{
		Code:    e.Value,
		Message: e.Message,
		Cause:   e.marshalCause(),
	})
}

func (e errorImpl) marshalCause() json.RawMessage {
	if e.Cause == nil {
		return nil
	}

	var out []byte

	// Voluntarily ignoring the marshalling errors as there's nothing we
	// can do about it.
	if impl, ok := e.Cause.(errorImpl); ok {
		out, _ = json.Marshal(impl)
	} else {
		out, _ = json.Marshal(e.Cause.Error())
	}

	return out
}

func determineCommonErrorMessage(code ErrorCode) string {
	switch code {
	case GenericErrorCode:
		return "An unexpected error occurred"
	case NotImplementedCode:
		return "Not implemented"
	default:
		return ""
	}
}
