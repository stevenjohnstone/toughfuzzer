package bigevent

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"hash/crc32"
	"strconv"
)

var strSum [sha256.Size]byte

func init() {
	const str = " is on! F$ck COVID"
	strSum = sha256.Sum256([]byte(str))
}

func match(data []byte) bool {
	sum := sha256.Sum256(data)
	return bytes.Compare(strSum[:], sum[:]) == 0
}

func bigEvent(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// collection of literal should handle this one
	if string(data[:4]) != "race" {
		return false
	}
	data = data[4:]

	if len(data) < 4 {
		return false
	}

	// We expect ascii encoded 2021...note that we're comparing against an int and sonar
	// will need to do some magic to make sure that an ascii string representation is tried
	if year, err := strconv.Atoi(string(data[:4])); err != nil || year != 2021 {
		return false
	}

	data = data[4:]

	if len(data) < 4 {
		return false
	}

	if !match(data[:len(data)-4]) {
		return false
	}

	// for fun, let's suppose that there's a CRC checksum too!
	sum := crc32.ChecksumIEEE(data[:len(data)-4])
	return binary.BigEndian.Uint32(data[len(data)-4:]) == sum
}

// FuzzBigEvent let's see if we can reach the end in record time!
func FuzzBigEvent(data []byte) int {
	if bigEvent(data) {
		panic("completed the obstacle course!")
	}
	return 0
}
