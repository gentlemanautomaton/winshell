// +build windows

package shobjidl

import (
	"syscall"
	"unsafe"

	"github.com/gentlemanautomaton/winshell/shellclass"
	"github.com/gentlemanautomaton/winshell/shellinterface"
	"github.com/go-ole/go-ole"
	"github.com/scjalliance/comutil"
)

// IShellLinkVtbl represents the component object model virtual
// function table for the IShellLinkW interface.
type IShellLinkVtbl struct {
	ole.IUnknownVtbl
	GetPath             uintptr
	GetIDList           uintptr
	SetIDList           uintptr
	GetDescription      uintptr
	SetDescription      uintptr
	GetWorkingDirectory uintptr
	SetWorkingDirectory uintptr
	GetArguments        uintptr
	SetArguments        uintptr
	GetHotkey           uintptr
	SetHotkey           uintptr
	GetShowCmd          uintptr
	SetShowCmd          uintptr
	GetIconLocation     uintptr
	SetIconLocation     uintptr
	SetRelativePath     uintptr
	Resolve             uintptr
	SetPath             uintptr
}

// IShellLink represents the component object model interface for
// manipulation of shell links. It implements the IShellLinkW interface.
type IShellLink struct {
	ole.IUnknown
}

// NewIShellLink creates a new instance of an IShellLink object.
func NewIShellLink() (*IShellLink, error) {
	iface, err := comutil.CreateObject(shellclass.ShellLink, shellinterface.ShellLink)
	return (*IShellLink)(unsafe.Pointer(iface)), err
}

// VTable returns the component object model virtual function table for the
// object.
func (v *IShellLink) VTable() *IShellLinkVtbl {
	return (*IShellLinkVtbl)(unsafe.Pointer(v.RawVTable))
}

// GetPath retrieves the path of the shell link.
//
// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-ishelllinkw-getpath
func (v *IShellLink) GetPath() (path string, err error) {
	const maxChars = 65535
	var buffer [maxChars]uint16
	hr, _, _ := syscall.Syscall6(
		v.VTable().GetPath,
		5,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&buffer[0])),
		maxChars,
		0,
		0,
		0)
	if hr != 0 {
		return "", ole.NewError(hr)
	}
	return syscall.UTF16ToString(buffer[:]), nil
}

// SetPath sets the path of the shell link.
//
// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-ishelllinkw-setpath
func (v *IShellLink) SetPath(path string) error {
	bpath, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	hr, _, _ := syscall.Syscall(
		v.VTable().SetPath,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(bpath)),
		0)
	if hr != 0 {
		return ole.NewError(hr)
	}
	return nil
}

// GetDescription retrieves the description of the shell link.
//
// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-ishelllinkw-getdescription
func (v *IShellLink) GetDescription() (description string, err error) {
	const maxChars = 65535
	var buffer [maxChars]uint16
	hr, _, _ := syscall.Syscall(
		v.VTable().GetPath,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(&buffer[0])),
		maxChars)
	if hr != 0 {
		return "", ole.NewError(hr)
	}
	return syscall.UTF16ToString(buffer[:]), nil
}

// SetDescription sets the description of the shell link.
//
// https://docs.microsoft.com/en-us/windows/win32/api/shobjidl_core/nf-shobjidl_core-ishelllinkw-setdescription
func (v *IShellLink) SetDescription(description string) error {
	bdescription, err := syscall.UTF16PtrFromString(description)
	if err != nil {
		return err
	}
	hr, _, _ := syscall.Syscall(
		uintptr(v.VTable().SetPath),
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(bdescription)),
		0)
	if hr != 0 {
		return ole.NewError(hr)
	}
	return nil
}

// PersistStream returns an implementation of the IPersistStream interface,
// which can be used to retrieve a binary representation of the shell link.
//
// It is the callers responsibility to call Release on the returned object
// when finished with it.
func (v *IShellLink) PersistStream() (*IPersistStream, error) {
	idispatch, err := v.QueryInterface(comutil.GUID(shellinterface.PersistStream))
	if err != nil {
		return nil, err
	}
	return (*IPersistStream)(unsafe.Pointer(idispatch)), nil
}
