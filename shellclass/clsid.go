// Package shellclass provides component object model class
// identifiers for the Windows shell.
package shellclass

import (
	"github.com/google/uuid"
)

var (
	// ShellLink is the component object model identifier of the ShellLink
	// class (CLSID_ShellLink).
	//
	//	{00021401-0000-0000-C000-000000000046}
	ShellLink = uuid.UUID{0x00, 0x02, 0x14, 0x01, 0x00, 0x00, 0x00, 0x00, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}

	// MyComputer is the component object model identifier of the MyComputer
	// class (CLSID_MyComputer).
	//
	//	{20D04FE0-3AEA-1069-A2D8-08002B30309D}
	MyComputer = uuid.UUID{0x20, 0xD0, 0x4F, 0xE0, 0x3A, 0xEA, 0x10, 0x69, 0xA2, 0xD8, 0x08, 0x00, 0x2B, 0x30, 0x30, 0x9D}
)
