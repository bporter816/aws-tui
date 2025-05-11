package template

import (
	"encoding/hex"
)

func FormatSerial(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	}
	dst := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(dst, bytes)
	final := make([]byte, 3*len(bytes))
	for i := 0; i < len(bytes); i++ {
		final[3*i] = dst[2*i]
		final[3*i+1] = dst[2*i+1]
		final[3*i+2] = ':'
	}
	return string(final[:len(final)-1])
}

func Chunk(str string, length int) []string {
	chunks := make([]string, 0)
	for i := 0; i <= len(str)-length; i += length {
		chunks = append(chunks, str[i:i+length])
	}
	if len(str)%length != 0 {
		chunks = append(chunks, str[len(str)-len(str)%length:])
	}
	return chunks
}
