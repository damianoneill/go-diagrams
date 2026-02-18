package rand

import "math/rand/v2"

const charset = "abcdefghijlkmnopqrstuvwxyz"

func String(length int) string {
	return StringWithCharset(length, charset)
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}

	return string(b)
}
