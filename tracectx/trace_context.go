package tracectx

import (
	"crypto/rand"
	"io"
)

const Version byte = 0

func NewParent() (Parent, error) {
	var buf [24]byte
	if _, err := io.ReadFull(rand.Reader, buf[:]); err != nil {
		return Parent{}, err
	}

	parent := Parent{
		version: Version,
	}

	copy(parent.traceID[:], buf[:16])
	copy(parent.parentID[:], buf[17:])

	return parent, nil
}

type Parent struct {
	version  byte
	traceID  [16]byte
	parentID [8]byte
	flags    byte
}
