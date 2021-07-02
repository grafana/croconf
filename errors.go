package croconf

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrorMissing = errors.New("field is missing in config source") //TODO: improve

type BindFieldMissingError struct {
	SourceName string
	Field      string // search field
}

func NewBindFieldMissingError(srcName string, field string) *BindFieldMissingError {
	return &BindFieldMissingError{SourceName: srcName, Field: field}
}

func (e *BindFieldMissingError) Error() string {
	return fmt.Sprintf("field %s is missing in config source %s", e.Field, e.SourceName)
}

type BindValueError struct {
	Func  string // the failing function (like BindIntValueTo)
	Input string // the input
	Err   error  // the reason the conversion failed
}

func NewBindValueError(f string, input string, err error) *BindValueError {
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		err = numErr.Err
	}
	return &BindValueError{
		Func: f, Input: input, Err: err,
	}
}

// Error implements error interface
func (e *BindValueError) Error() string {
	return e.Func + ": " + "parsing " + strconv.Quote(e.Input) + ": " + e.Err.Error()
}

func (e *BindValueError) Unwrap() error { return e.Err }

func (e *BindValueError) withFuncName(funcName string) *BindValueError {
	return NewBindValueError(funcName, e.Input, e.Err)
}

type JSONSourceInitError struct {
	Data []byte // the failing data input
	Err  error
}

func NewJSONSourceInitError(data []byte, err error) *JSONSourceInitError {
	return &JSONSourceInitError{Data: data, Err: err}
}

// Error implements error interface
func (e *JSONSourceInitError) Error() string {
	return "source json initialization failed: data=" + string(e.Data)
}

func (e *JSONSourceInitError) Unwrap() error { return e.Err }
