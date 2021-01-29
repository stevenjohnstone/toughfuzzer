package sonar

const leet = "1337 string abcdefg"

func reverse(s string) string {
	r := make([]rune, len(s))
	for i, c := range s {
		r[len(s)-1-i] = c
	}
	return string(r)
}

func findString(s, t string) {
	if s == reverse(t) {
		panic("found string")
	}
}

// FuzzString challenges the fuzzer to find a simple string which is reversed
func FuzzString(data []byte) int {
	// note that if we swap the arguments of findString,
	// sonar will not help us find a matching string
	findString(string(data), leet)
	return 0
}
