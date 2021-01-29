package sonar

import (
	"encoding/binary"
	"hash/crc32"
)

// checksumOK returns true if the first four bytes of data
// are the CRC checksum of the remaining bytes
func checksumOK(data []byte) bool {
	if len(data) <= 4 {
		return false
	}
	sum := crc32.ChecksumIEEE(data[4:])
	return binary.BigEndian.Uint32(data[:4]) == sum
}

// FuzzCheckSum shows how sonar can modify inputs to have correct checksums
func FuzzCheckSum(data []byte) int {
	if checksumOK(data) {
		panic("found correct checksum")
	}
	return 0
}
