package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

func ShaHash() string {

	randomString := randomString(40)
	sha512Sum := sha512.Sum512_256([]byte(randomString))
	return hex.EncodeToString(sha512Sum[:])

}

func randomString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
