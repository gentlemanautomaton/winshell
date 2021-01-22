// +build windows

package shelllink_test

import (
	"fmt"

	"github.com/gentlemanautomaton/winshell/shelllink"
)

func ExamplePath() {
	data, err := shelllink.Path(`C:\Users`).MarshalBinary()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%x", data[:24])

	// Output: 4c0000000114020000000000c00000000000004683000000
}
