package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GetRandomShortLink() string {
	b := make([]byte, 6)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
