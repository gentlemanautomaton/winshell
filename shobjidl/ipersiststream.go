// +build windows

package shobjidl

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

// IPersistStreamVtbl represents the component object model virtual
// function table for the IPersistStream interface.
type IPersistStreamVtbl struct {
	IPersistVtbl
	IsDirty    uintptr
	Load       uintptr
	Save       uintptr
	GetSizeMax uintptr
}

// IPersistStream represents the component object model interface for
// marshaling binary representations of objects.
type IPersistStream struct {
	IPersist
}

// VTable returns the component object model virtual function table for the
// object.
func (v *IPersistStream) VTable() *IPersistStreamVtbl {
	return (*IPersistStreamVtbl)(unsafe.Pointer(v.RawVTable))
}

// IsDirty returns true if the object holds unsaved data.
//
// https://docs.microsoft.com/en-us/windows/win32/api/objidl/nf-objidl-ipersiststream-isdirty
func (v *IPersistStream) IsDirty() (bool, error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().IsDirty,
		1,
		uintptr(unsafe.Pointer(v)),
		0,
		0)
	switch hr {
	case 0:
		return true, nil
	case 1:
		return false, nil
	default:
		return true, ole.NewError(hr)
	}
}

// Load copies data from stream into v.
//
// https://docs.microsoft.com/en-us/windows/win32/api/objidl/nf-objidl-ipersiststream-load
func (v *IPersistStream) Load(stream *IStream) error {
	hr, _, _ := syscall.Syscall(
		v.VTable().Load,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(stream)),
		0)
	if hr != 0 {
		return ole.NewError(hr)
	}
	return nil
}

// Save copies data from v into stream.
//
// https://docs.microsoft.com/en-us/windows/win32/api/objidl/nf-objidl-ipersiststream-save
func (v *IPersistStream) Save(stream *IStream) error {
	hr, _, _ := syscall.Syscall(
		v.VTable().Save,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(stream)),
		0)
	if hr != 0 {
		return ole.NewError(hr)
	}
	return nil
}

// GetSizeMax returns the number of bytes needed to store a copy of the
// object.
//
// https://docs.microsoft.com/en-us/windows/win32/api/objidl/nf-objidl-ipersiststream-getsizemax
func (v *IPersistStream) GetSizeMax() (size uint64, err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().GetSizeMax,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&size)),
		0)
	if hr != 0 {
		return 0, ole.NewError(hr)
	}
	return
}
