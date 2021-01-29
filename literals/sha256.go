package literals

import (
	"bytes"
	"crypto/sha256"
)

var strSum [sha256.Size]byte

func init() {
	const str = "really too long to be guessed"
	strSum = sha256.Sum256([]byte(str))
}

func match(data []byte) bool {
	sum := sha256.Sum256(data)
	return bytes.Compare(strSum[:], sum[:]) == 0
}

// FuzzLiteral shows how a string literal will be extracted from the program to
// be used as an fuzzing input.
func FuzzLiteral(data []byte) int {
	if match(data) {
		panic("found string")
	}
	return 0
}
