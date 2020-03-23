package socketcore

import (
	"sync"
)

const (
	pageSize = 2048
)

var (
	memorypooltmp = &sync.Pool{
		New: func() interface{} {
			return make([]byte, pageSize)
		},
	}
)

// Alloc ...
func Alloc() []byte {
	return memorypooltmp.Get().([]byte)
}

// Free ..
func Free(b []byte) {
	if b != nil {
		memorypooltmp.Put(b)
	}
}
