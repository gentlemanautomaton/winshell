// +build windows

package shobjidl

import (
	"io"
	"sync"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

const (
	efail       = 0x80004005
	einvalidarg = 0x80070057
	enotimpl    = 0x80004001
)

// This virtual function table relies on Go "method expressions" to map
// pointer receivers to callbacks compatible with COM calling conventions.
//
// Each callback is created from a method expression (*T).Method, which
// is an alternative form of Method that takes the receiver (*T) as the first
// argument. This matches COM calling conventions exactly, which pass a
// pointer to the object as the first argument.
//
// When COM calls one of these virtual functions, the method is invoked
// exactly as it would be from Go code as T.Method().
//
// See: https://golang.org/ref/spec#Method_expressions
var streamBufferVTable = IStreamVtbl{
	ISequentialStreamVtbl: ISequentialStreamVtbl{
		IUnknownVtbl: ole.IUnknownVtbl{
			QueryInterface: syscall.NewCallback((*StreamBuffer).queryInterface),
			AddRef:         syscall.NewCallback((*StreamBuffer).addRef),
			Release:        syscall.NewCallback((*StreamBuffer).release),
		},
		Read:  syscall.NewCallback((*StreamBuffer).read),
		Write: syscall.NewCallback((*StreamBuffer).write),
	},
	Seek:         syscall.NewCallback((*StreamBuffer).seek),
	SetSize:      syscall.NewCallback((*StreamBuffer).setSize),
	CopyTo:       syscall.NewCallback((*StreamBuffer).copyTo),
	Commit:       syscall.NewCallback((*StreamBuffer).commit),
	Revert:       syscall.NewCallback((*StreamBuffer).revert),
	LockRegion:   syscall.NewCallback((*StreamBuffer).lockRegion),
	UnlockRegion: syscall.NewCallback((*StreamBuffer).unlockRegion),
	Stat:         syscall.NewCallback((*StreamBuffer).stat),
	Clone:        syscall.NewCallback((*StreamBuffer).clone),
}

// Raymond Chen has a wonderfully succinct blog post about COM object layouts:
// https://devblogs.microsoft.com/oldnewthing/20040205-00/?p=40733

// StreamBuffer is an in-memory implementation of a COM IStream.
type StreamBuffer struct {
	vtable *IStreamVtbl
	mutex  sync.RWMutex
	data   []byte
	offset int
	log    []string
}

// NewStreamBuffer returns a stream buffer that implement IStream.
func NewStreamBuffer() *StreamBuffer {
	return &StreamBuffer{
		vtable: &streamBufferVTable,
	}
}

// IStream returns a component object model representation of the stream.
func (b *StreamBuffer) IStream() *IStream {
	return (*IStream)(unsafe.Pointer(b))
}

// Bytes returns the content of the stream buffer.
func (b *StreamBuffer) Bytes() []byte {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.data
}

func (b *StreamBuffer) queryInterface() uintptr {
	return enotimpl
}

func (b *StreamBuffer) addRef() uintptr {
	return 0
}

func (b *StreamBuffer) release() uintptr {
	return 0
}

// read copies bytes from b to buf. The buf pointer is interpreted as a
// pointer to a byte array of size bufSize. If n is not nil, The number of
// bytes copied is stored in it.
func (b *StreamBuffer) read(buf *byte, bufSize uint32, n *uint32) uintptr {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Return early if the offset is already at the end of the buffer
	if b.offset >= len(b.data) {
		if n != nil {
			*n = 0
		}
		return 0
	}

	// Transmute the output buffer pointer into a byte slice
	const maxChunk = 4096
	if bufSize > maxChunk {
		bufSize = maxChunk
	}
	out := (*[maxChunk]byte)(unsafe.Pointer(buf))[:bufSize:bufSize]

	// Copy bytes to the output buffer
	_n := copy(out, b.data[b.offset:])

	// Advance the offset
	b.offset += _n

	// Report the number of bytes read, if requested
	if n != nil {
		*n = uint32(_n)
	}

	return 0
}

// write copies bytes from buf to b. The buf pointer is interpreted as a
// pointer to a byte array of size bufSize. If n is not nil, The number of
// bytes copied is stored in it.
func (b *StreamBuffer) write(buf *byte, bufSize uint32, n *uint32) uintptr {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Transmute the input buffer pointer into a byte slice
	const maxChunk = 4096
	if bufSize > maxChunk {
		bufSize = maxChunk
	}
	in := (*[maxChunk]byte)(unsafe.Pointer(buf))[:bufSize:bufSize]

	// Grow the local buffer if necessary
	if size, needed := len(b.data), b.offset+len(in); size < needed {
		b.data = append(b.data, make([]byte, needed-size)...)
	}

	// Copy bytes from the input buffer
	_n := copy(b.data[b.offset:], in)

	// Advance the offset
	b.offset += _n

	// Report the number of bytes written, if requested
	if n != nil {
		*n = uint32(_n)
	}

	return 0
}

// seek adjust the position of the current offset within the buffer.
func (b *StreamBuffer) seek(offset int64, whence uint32, pos *uint64) uintptr {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var off int
	switch whence {
	case io.SeekStart:
		off = int(offset)
	case io.SeekCurrent:
		off = b.offset + int(offset)
	case io.SeekEnd:
		off = len(b.data) + int(offset)
	default:
		return einvalidarg
	}

	if off < 0 {
		return einvalidarg
	}

	b.offset = off

	if pos != nil {
		*pos = uint64(b.offset)
	}

	return 0
}

// setSize ensure that the local buffer is at least size bytes long.
func (b *StreamBuffer) setSize(size int64) uintptr {
	if size, needed := len(b.data), int(size); size < needed {
		b.data = append(b.data, make([]byte, needed-size)...)
	}
	return 0
}

func (b *StreamBuffer) copyTo() uintptr {
	return enotimpl
}

func (b *StreamBuffer) commit() uintptr {
	return enotimpl
}

func (b *StreamBuffer) revert() uintptr {
	return enotimpl
}

func (b *StreamBuffer) lockRegion() uintptr {
	return enotimpl
}

func (b *StreamBuffer) unlockRegion() uintptr {
	return enotimpl
}

func (b *StreamBuffer) stat() uintptr {
	return enotimpl
}

func (b *StreamBuffer) clone() uintptr {
	return enotimpl
}
