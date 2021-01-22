// +build windows

package shobjidl

import (
	"unsafe"
)

// IStreamVtbl represents the component object model virtual
// function table for the IStream interface.
type IStreamVtbl struct {
	ISequentialStreamVtbl
	Seek         uintptr
	SetSize      uintptr
	CopyTo       uintptr
	Commit       uintptr
	Revert       uintptr
	LockRegion   uintptr
	UnlockRegion uintptr
	Stat         uintptr
	Clone        uintptr
}

// IStream represents the component object model interface for
// marshaling binary representations of objects.
type IStream struct {
	ISequentialStream
}

// VTable returns the component object model virtual function table for the
// object.
func (v *IStream) VTable() *IStreamVtbl {
	return (*IStreamVtbl)(unsafe.Pointer(v.RawVTable))
}
