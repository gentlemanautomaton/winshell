package shobjidl_test

import (
	"fmt"

	"github.com/gentlemanautomaton/winshell/shobjidl"
	"github.com/scjalliance/comshim"
)

func ExampleIShellLink() {
	// Make sure COM is initialized for the duration of the work
	comshim.Add(1)
	defer comshim.Done()

	// Prepare a new IShellLink COM instance
	link, err := shobjidl.NewIShellLink()
	if err != nil {
		panic(err)
	}
	defer link.Release()

	// Set the path used by the shortcut
	if err := link.SetPath(`C:\Users`); err != nil {
		panic(err)
	}

	// Ask the link for an IPersistStream COM interface that we can use
	// to serialize its data
	stream, err := link.PersistStream()
	if err != nil {
		panic(err)
	}
	defer stream.Release()

	// Write the link's data to an in-memory buffer
	buf := shobjidl.NewStreamBuffer()
	if err := stream.Save(buf.IStream()); err != nil {
		panic(err)
	}

	// Write the link data to a file
	//ioutil.WriteFile("test.lnk", buf.Bytes(), 0644)

	fmt.Printf("Success")

	// Output: Success
}
