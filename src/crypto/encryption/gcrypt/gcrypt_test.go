package gcrypt_test

import (
	"drm-blockchain/src/crypto/encryption/gcrypt"
	"drm-blockchain/src/crypto/random"
	ez "drm-blockchain/src/utils/test"
	"testing"
)

func Test__GCrypto__EncryptThenDecrypt__Should__ReturnOriginalValue(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("This is a key")

	g := gcrypt.New(128)
	encrypted := g.Encrypt(data, key)
	decrypted := g.Decrypt(encrypted, key)

	if string(decrypted) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(decrypted))
	}
}

func Test__GCrypto__ToBytesThenFromBytes__Should__ReturnOriginalValue(t *testing.T) {
	data := []byte("Hello, World!")

	g := gcrypt.New(128)
	encoded := g.EncodeToBytes(data)
	decoded := g.Decode(g.FromBytes(encoded))

	if string(decoded) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(decoded))
	}
}

func Test__GCrypto__ToBytesThenFromBytes__Should__ReturnOriginalValue2(t *testing.T) {
	ez := ez.New(t)
	crypt := gcrypt.New(64)

	rnd, _ := random.GenerateBytes(1)
	data := crypt.EncodeToBytes(rnd)
	encData := crypt.Encode(rnd)
	expData := crypt.FromBytes(data)

	ez.AssertAreEqual(expData, encData)
}
