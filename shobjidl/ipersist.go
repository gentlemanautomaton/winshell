// +build windows

package shobjidl

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/google/uuid"
)

// IPersistVtbl represents the component object model virtual
// function table for the IPersist interface.
type IPersistVtbl struct {
	ole.IUnknownVtbl
	GetClassID uintptr
}

// IPersist represents the component object model interface for
// peristable objects.
type IPersist struct {
	ole.IUnknown
}

// VTable returns the component object model virtual function table for the
// object.
func (v *IPersist) VTable() *IPersistVtbl {
	return (*IPersistVtbl)(unsafe.Pointer(v.RawVTable))
}

// GetClassID retrieves the class identifier of the class that is capable of
// manipulating the object's data.
//
// https://docs.microsoft.com/en-us/windows/win32/api/objidl/nf-objidl-ipersist-getclassid
func (v *IPersist) GetClassID() (classID uuid.UUID, err error) {
	hr, _, _ := syscall.Syscall(
		uintptr(v.VTable().GetClassID),
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&classID[0])),
		0)
	if hr != 0 {
		err = ole.NewError(hr)
	}
	return
}
