package croconf

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrorMissing = errors.New("field is missing in config source") //TODO: improve

type BindNameMissingError struct {
	Func string // the failing function (like BindIntValueTo)
	Name string // search field
}

func NewBindNameMissingError(f string, name string) *BindNameMissingError {
	return &BindNameMissingError{Func: f, Name: name}
}

func (e *BindNameMissingError) Error() string {
	return fmt.Sprintf("%s: binding name %s not found in config source", e.Func, e.Name)
}

type BindValueError struct {
	Func  string // the failing function (like BindIntValueTo)
	Input string // the input
	Err   error  // the reason the conversion failed
}

func NewBindValueError(f string, input string, err error) *BindValueError {
	if numErr, ok := err.(*strconv.NumError); ok {
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

func (e *BindValueError) withFuncName(funcName string) *BindValueError {
	return NewBindValueError(funcName, e.Input, e.Err)
}

// setErrorFunc changes func name in the custom error
func setErrorFunc(src error, funcName string) {
	bindErr, ok := src.(*BindValueError)
	if !ok {
		// do nothing if unexpected error type
		return
	}
	bindErr.Func = funcName
}
