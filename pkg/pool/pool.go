package pool

import (
	"bytes"
	"sync"
)

// ByteBufferPool provides reusable byte buffers to reduce allocations
var ByteBufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GetBuffer retrieves a buffer from the pool
func GetBuffer() *bytes.Buffer {
	buf := ByteBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *bytes.Buffer) {
	if buf != nil {
		buf.Reset()
		ByteBufferPool.Put(buf)
	}
}

// SlicePool for []interface{} slices used in responses
var SlicePool = sync.Pool{
	New: func() interface{} {
		// Pre-allocate with reasonable capacity
		s := make([]interface{}, 0, 20)
		return &s
	},
}

// GetSlice retrieves a slice from the pool
func GetSlice() *[]interface{} {
	s := SlicePool.Get().(*[]interface{})
	*s = (*s)[:0] // Reset length but keep capacity
	return s
}

// PutSlice returns a slice to the pool
func PutSlice(s *[]interface{}) {
	if s != nil {
		*s = (*s)[:0]
		SlicePool.Put(s)
	}
}

// StringBuilderPool for string building operations
var StringBuilderPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// MapPool for map[string]interface{} used in JSON responses
var MapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{}, 8)
	},
}

// GetMap retrieves a map from the pool
func GetMap() map[string]interface{} {
	m := MapPool.Get().(map[string]interface{})
	// Clear the map
	for k := range m {
		delete(m, k)
	}
	return m
}

// PutMap returns a map to the pool
func PutMap(m map[string]interface{}) {
	if m != nil {
		for k := range m {
			delete(m, k)
		}
		MapPool.Put(m)
	}
}
