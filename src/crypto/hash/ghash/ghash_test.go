package ghash_test

import (
	"drm-blockchain/src/crypto/encryption/gcrypto"
	"drm-blockchain/src/crypto/hash/ghash"
	ez "drm-blockchain/src/utils/test"
	"testing"
)

func Test__GHashOverGCrypto(t *testing.T) {
	data := []byte("Hello, World!")
	key := []byte("This is a key")

	crypt := gcrypto.New(128)

	dataHash := ghash.New(500, 128, 128, key)
	dataHash.SetNonce([]byte("dataNonce"))
	dataHash.Add(crypt.EncodeToBytes(data))
	dataHashBytes := dataHash.GetDigest()

	keyHash := ghash.New(500, 128, 128, key)
	keyNonce := []byte("keyNonce")
	keyHash.SetNonce(keyNonce)
	encodedKey := crypt.EncodeKeyToBytes(key, len(data))
	keyHash.Add(encodedKey)

	encrypted := crypt.Encrypt(data, key)
	encryptedNonceState := dataHash.GetNonceState()
	keyNonceState := keyHash.GetNonceState()

	for i := range encryptedNonceState {
		encryptedNonceState[i] = encryptedNonceState[i].Add(keyNonceState[i])
	}

	encryptedHash := ghash.New(500, 128, 128, key)
	encryptedHash.SetNonceState(encryptedNonceState)
	encryptedHash.Add(encrypted)

	encryptedHash.Remove(encodedKey)
	encryptedHash.RemoveNonce(keyNonce)
	encryptedHashBytes := encryptedHash.GetDigest()

	ez := ez.New(t)
	ez.AssertAreEqual(dataHashBytes, encryptedHashBytes)
}
