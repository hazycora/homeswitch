package crypto

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomString(byteLength int64) (str string, err error) {
	b := make([]byte, byteLength)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	str = base64.RawURLEncoding.EncodeToString(b)
	return
}
