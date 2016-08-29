// tracer.go
package tracer

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

type nilTracer struct{}

type tracer struct {
	out io.Writer
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
func (t *nilTracer) Trace(a ...interface{}) {

}
func Off() Tracer {
	return &nilTracer{}
}
