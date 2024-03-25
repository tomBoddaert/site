package main

import (
	"strings"
)

type MaxStringBuilder struct {
	MaxLen uint
	strings.Builder
}

func (b *MaxStringBuilder) Write(p []byte) (int, error) {
	if b.Builder.Len() < int(b.MaxLen) {
		writing := min(len(p), int(b.MaxLen)-b.Builder.Len())
		return b.Builder.Write(p[:writing])
	}

	return len(p), nil
}

func (b *MaxStringBuilder) Len() int {
	return b.Builder.Len()
}

func (b *MaxStringBuilder) String() string {
	return b.Builder.String()
}
