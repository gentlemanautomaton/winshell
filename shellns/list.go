package shellns

import (
	"encoding/binary"
	"fmt"
)

// List is a list of shell namespace item IDs.
type List []Item

// Size returns the number of bytes needed to marshal the list.
func (list List) Size() int {
	var total int

	// Add two bytes for the 16-bit list header
	total += 2

	for i := range list {
		// Add two bytes for the 16-bit item header
		total += 2
		// Add the size of the item itself
		total += len(list[i])
	}

	// Add two bytes for the list terminal
	total += 2

	return total
}

// MarshalBinaryTo writes a binary representation of the shell link item ID
// list to data.
//
// The list will be marshaled as a LinkTargetIDList:
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-shllink/881d7a83-07a5-4702-93e3-f9fc34c3e1e4
func (list List) MarshalBinaryTo(data []byte) (err error) {
	// Calculate the number of bytes needed
	size := list.Size()

	// Make sure we aren't going to overflow a 16 bit integer in any of
	// our conversions
	if size > 65535 {
		return fmt.Errorf("the shell namespace item ID list requires %d bytes, which exceeds the limit of 65535", size)
	}

	// Make sure the buffer is sufficiently sized.
	if len(data) < size {
		return fmt.Errorf("the shell namespace item ID list requires %d bytes, but the buffer provided holds %d bytes", size, len(data))
	}

	// Write a 16 bit list header with the size
	binary.LittleEndian.PutUint16(data[0:2], uint16(len(list)))
	offset := 2

	for _, item := range list {
		// Write a 16 bit item header with the size
		binary.LittleEndian.PutUint16(data[offset:offset+2], uint16(len(item)))
		offset += 2

		// Write the item bytes
		copy(data[offset:], item)
		offset += len(item)
	}

	// Write a 16 bit list terminal
	binary.LittleEndian.PutUint16(data[offset:offset+2], 0)

	return nil
}

// MarshalBinary returns a binary representation of the shell link item
// ID list.
//
// The list will be marshaled as a LinkTargetIDList:
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-shllink/881d7a83-07a5-4702-93e3-f9fc34c3e1e4
func (list List) MarshalBinary() (data []byte, err error) {
	data = make([]byte, list.Size())
	return data, list.MarshalBinaryTo(data)

}
