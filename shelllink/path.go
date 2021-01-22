// +build windows

package shelllink

import (
	"fmt"

	"github.com/gentlemanautomaton/winshell/shobjidl"
	"github.com/scjalliance/comshim"
)

// Path is a file system path that can be marshaled as a shell link.
type Path string

// MarshalBinary returns a binary representation of a shell link for the path.
func (p Path) MarshalBinary() ([]byte, error) {
	// Make sure COM is initialized for the duration of the work
	comshim.Add(1)
	defer comshim.Done()

	// Prepare a new IShellLink COM instance
	link, err := shobjidl.NewIShellLink()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare IShellLink COM interface: %v", err)
	}
	defer link.Release()

	// Set the path for the shortcut
	if err := link.SetPath(string(p)); err != nil {
		return nil, fmt.Errorf("failed to set shell link path: %v", err)
	}

	// Ask the link for an IPersistStream COM interface that we can use
	// to serialize its data
	stream, err := link.PersistStream()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare IPersistStream COM interface: %v", err)
	}
	defer stream.Release()

	// Write the link's data to an in-memory buffer
	buf := shobjidl.NewStreamBuffer()
	if err := stream.Save(buf.IStream()); err != nil {
		return nil, fmt.Errorf("failed to save shell link data to stream buffer: %v", err)
	}

	return buf.Bytes(), nil
}
