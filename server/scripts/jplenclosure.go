package scripts

import (
	"github.com/jplorg/jpl/go/v2/jpl"
	"github.com/jplorg/jpl/go/v2/library"
)

func enclose(fn func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error)) jpl.JPLFunc {
	return &enclosure{fn: fn}
}

type enclosure struct {
	fn func(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, input any, args ...any) ([]any, error)
}

func (e *enclosure) Call(runtime jpl.JPLRuntime, signal jpl.JPLRuntimeSignal, next jpl.JPLPiper, input any, args ...any) ([]any, error) {
	results, err := e.fn(runtime, signal, input, args...)
	if err != nil {
		return nil, err
	}

	return library.MuxAll([][]any{results}, library.NewPiperMuxer(next))
}

func (e *enclosure) IsSame(other jpl.JPLFunc) bool {
	return e == other
}
