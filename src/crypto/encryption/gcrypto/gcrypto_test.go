package gcrypto_test

import (
	"drm-blockchain/src/crypto/encryption/gcrypto"
	"testing"
)

func Test__GCrypto__EncryptThenDecrypt__Should__ReturnOriginalValue(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("This is a key")

	g := gcrypto.New(128)
	encrypted := g.Encrypt(data, key)
	decrypted := g.Decrypt(encrypted, key)

	if string(decrypted) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(decrypted))
	}
}
