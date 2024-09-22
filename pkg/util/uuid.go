package util

import (
	"crypto/rand"
	"fmt"
)

func GenerateUUIDv4() string {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		return ""
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4 (random)
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10xxxxxx

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:])
}
