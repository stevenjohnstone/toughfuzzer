// +build gofuzz

package sonar

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
)

func findInt(i int) {
	if 1337 == i {
		panic("Found int")
	}
}

// FuzzIntegerBigEndian shows how sonar can discover integers to use as inputs to increase coverage
func FuzzIntegerBigEndian(data []byte) int {
	// interpret data as big-endian encoded
	findInt(int(binary.BigEndian.Uint64(data)))
	return 0
}

// FuzzIntegerLittleEndian shows how sonar can discover integers to use as inputs to increase coverage
func FuzzIntegerLittleEndian(data []byte) int {
	// interpret data as little-endian encoded
	findInt(int(binary.LittleEndian.Uint64(data)))
	return 0
}

// FuzzIntegerDecimalString interprets input as a decimal string
func FuzzIntegerDecimalString(data []byte) int {
	i, err := strconv.Atoi(string(data))
	if err != nil {
		return 0
	}
	findInt(i)
	return 1
}

// FuzzIntegerHexString interprets input as a hexadecimal string
func FuzzIntegerHexString(data []byte) int {
	i, err := hex.DecodeString(string(data))
	if err != nil {
		return 0
	}
	FuzzIntegerBigEndian(i)
	return 1
}
