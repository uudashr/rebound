package rebound

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// EventHandle is a function type that handle the event.
// The function should return an error if the handling is failed.
// The function form is:
//
//		func(event Event) error
//
//	 where the Event is the event type (struct) that will be handled.
//
// Example:
//
//	eventually.HandleEvent(func(event OrderCompleted) error {
//		// handle the event
//		return nil
//	})
type EventHandler any

type Rebound struct {
	handlers map[string]EventHandler
	Decoder  Decoder
}

func (r *Rebound) ReactTo(eventName string, fn EventHandler) {
	if eventName == "" {
		panic("rebound: event name is empty")
	}

	err := ValidateHandler(fn)
	if err != nil {
		panic(err)
	}

	if r.handlers == nil {
		r.handlers = make(map[string]EventHandler)
	}

	_, exists := r.handlers[eventName]
	if exists {
		panic(fmt.Sprintf("rebound: event %q already has a handler", eventName))
	}

	r.handlers[eventName] = fn
}

func (r *Rebound) Dispatch(eventName string, data []byte) error {
	if eventName == "" {
		return fmt.Errorf("rebound: event name is empty")
	}

	fn := r.handlers[eventName]
	if fn == nil {
		return fmt.Errorf("rebound: no handler for event %q", eventName)
	}

	fnType := reflect.TypeOf(fn)
	event := reflect.New(fnType.In(0))
	err := r.decode(data, event.Interface())
	if err != nil {
		return fmt.Errorf("rebound: failed to unmarshal event data: %w", err)
	}

	fnValue := reflect.ValueOf(fn)
	retVals := fnValue.Call([]reflect.Value{event.Elem()})
	if !retVals[0].IsNil() {
		return retVals[0].Interface().(error)
	}

	return nil
}

func (r *Rebound) decode(data []byte, v interface{}) error {
	decoder := r.Decoder
	if decoder == nil {
		decoder = DefaultDecoder
	}

	return decoder.Decode(data, v)
}

func ValidateHandler(fn EventHandler) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("rebound: fn EventHandler is not a function (got: %v)", fnType.Kind())
	}

	if fnType.NumIn() != 1 {
		return fmt.Errorf("rebound: fn EventHandler should have 1 input parameter (got: %d)", fnType.NumIn())
	}

	if fnType.NumOut() != 1 {
		return fmt.Errorf("rebound: fn EventHandler should have 1 output parameter (got: %d)", fnType.NumOut())
	}

	if fnType.In(0).Kind() != reflect.Struct {
		return fmt.Errorf("rebound: fn EventHandler input parameter should be a struct (got: %v)", fnType.In(0).Kind())
	}

	if fnType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("rebound: fn EventHandler output parameter should be an error (got: %v)", fnType.Out(0))
	}

	return nil
}

type Decoder interface {
	Decode(data []byte, v interface{}) error
}

type DecodeFunc func(data []byte, v interface{}) error

func (f DecodeFunc) Decode(data []byte, v interface{}) error {
	return f(data, v)
}

var JsonDecoder = DecodeFunc(json.Unmarshal)
var DefaultDecoder = JsonDecoder
