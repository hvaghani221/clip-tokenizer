package main

import (
	"strings"

	"github.com/fatih/color"
)

type colorBuilder struct {
	builder     strings.Builder
	keyWriter   func(string, ...any) string
	valueWriter func(string, ...any) string
	keyWritten  bool
}

func NewColorBuilder() *colorBuilder {
	return &colorBuilder{
		keyWriter:   color.MagentaString,
		valueWriter: color.CyanString,
	}
}

func (cb *colorBuilder) WriteKeyValue(key string, value string, params ...any) {
	if cb.keyWritten {
		cb.builder.WriteString(", ")
	}
	cb.builder.WriteString(cb.keyWriter(key))
	cb.builder.WriteString(": ")
	cb.builder.WriteString(cb.valueWriter(value, params...))
	cb.keyWritten = true
}

func (cb *colorBuilder) WriteString(str string) {
	cb.builder.WriteString(str)
}

func (cb *colorBuilder) LineBreak() {
	cb.builder.WriteByte('\n')
	cb.keyWritten = false
}

func (cb *colorBuilder) Grow(n int) {
	cb.builder.Grow(n)
}

func (cb *colorBuilder) String() string {
	return cb.builder.String()
}
