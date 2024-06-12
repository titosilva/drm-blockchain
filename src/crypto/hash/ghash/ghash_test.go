package ghash_test

import (
	"drm-blockchain/src/crypto/encryption/gcrypt"
	"drm-blockchain/src/crypto/hash/ghash"
	"drm-blockchain/src/crypto/random"
	"drm-blockchain/src/math/uintp"
	ez "drm-blockchain/src/utils/test"
	"testing"
)

func Test__GHash__AddThenRemove__ShouldEqual__Original(t *testing.T) {
	ez := ez.New(t)

	data, _ := random.GenerateBytes(64)
	base, _ := random.GenerateBytes(64)

	hash := ghash.NewWithParams(1, 64, 128, nil)
	hash.AddBytes(base)
	hash.AddBytes(data)
	hash.RemoveBytes(data)

	baseHash := ghash.NewWithParams(1, 64, 128, nil)
	baseHash.AddBytes(base)

	ez.AssertAreEqual(hash.GetDigest(), baseHash.GetDigest())
}

func Test__GHash__AddThenRemove__ShouldEqual__Original__WithGCryptoMethodsWithoutByteConversion(t *testing.T) {
	ez := ez.New(t)
	crypt := gcrypt.New(64)

	rnd, _ := random.GenerateBytes(1)
	encData := crypt.Encode(rnd)

	rnd, _ = random.GenerateBytes(1)
	encKey := crypt.ExpandKey(rnd, len(encData))

	hash := ghash.NewWithParams(1, 64, 128, nil)
	hash.AddBlocks(encKey)
	hash.AddBlocks(encData)

	encrypted := make([]*uintp.UintP, len(encData))
	for i := 0; i < len(encData); i++ {
		encrypted[i] = uintp.Clone(encData[i]).Add(encKey[i])
	}

	encHash := ghash.NewWithParams(1, 64, 128, nil)
	encHash.AddBlocks(encrypted)

	ez.AssertAreEqual(hash.GetDigest(), encHash.GetDigest())
}

func Test__GHash__AddBytesAndAddBlocks__ShouldEqual(t *testing.T) {
	ez := ez.New(t)
	crypt := gcrypt.New(64)

	rnd, _ := random.GenerateBytes(1)
	data := crypt.EncodeToBytes(rnd)
	encData := crypt.Encode(rnd)
	expData := crypt.FromBytes(data)
	ez.AssertAreEqual(expData, encData)

	rnd, _ = random.GenerateBytes(1)
	key := crypt.ExpandKeyToBytes(rnd, len(encData))
	encKey := crypt.ExpandKey(rnd, len(encData))
	expKey := crypt.FromBytes(key)
	ez.AssertAreEqual(expKey, encKey)

	hash := ghash.NewWithParams(1, 64, 128, nil)
	hash.AddBytes(key)
	hash.AddBytes(data)

	blockHash := ghash.NewWithParams(1, 64, 128, nil)
	blockHash.AddBlocks(encKey)
	blockHash.AddBlocks(encData)

	ez.AssertAreEqual(hash.GetDigest(), blockHash.GetDigest())
}

func Test__GHash__AddThenRemove__ShouldEqual__Original__WithGCryptoMethodsWithByteConversions(t *testing.T) {
	ez := ez.New(t)
	crypt := gcrypt.New(64)

	rnd, _ := random.GenerateBytes(1)
	data := crypt.EncodeToBytes(rnd)
	encData := crypt.Encode(rnd)
	expData := crypt.FromBytes(data)
	ez.AssertAreEqual(expData, encData)

	rnd, _ = random.GenerateBytes(1)
	key := crypt.ExpandKeyToBytes(rnd, len(encData))
	encKey := crypt.ExpandKey(rnd, len(encData))
	expKey := crypt.FromBytes(key)
	ez.AssertAreEqual(expKey, encKey)

	hash := ghash.NewWithParams(1, 64, 128, nil)
	hash.AddBytes(key)
	hash.AddBytes(data)

	encrypted := make([]*uintp.UintP, len(encData))
	for i := 0; i < len(encData); i++ {
		encrypted[i] = uintp.Clone(encData[i]).Add(encKey[i])
	}

	encHash := ghash.NewWithParams(1, 64, 128, nil)
	encHash.AddBlocks(encrypted)

	ez.AssertAreEqual(hash.GetDigest(), encHash.GetDigest())
}

func Test__GHashOverGCrypto(t *testing.T) {
	ez := ez.New(t)
	data := []byte("Hello, World!")
	key := []byte("This is a key")

	crypt := gcrypt.New(64)

	dataHash := ghash.NewWithParams(1, 64, 128, nil)
	encodedData := crypt.Encode(data)
	encodedDataBytes := crypt.EncodeToBytes(data)
	dataHash.AddBytes(encodedDataBytes)
	dataDigest := dataHash.GetDigest()

	encrypted := crypt.Encrypt(data, key)

	encryptedHash := ghash.NewWithParams(1, 64, 128, nil)
	encryptedHash.AddBytes(encrypted)
	encodedKey := crypt.ExpandKeyToBytes(key, len(encodedData))
	encryptedHash.RemoveBytes(encodedKey)

	recoveredHash := encryptedHash.GetDigest()

	ez.AssertAreEqual(crypt.Decrypt(encrypted, key), data)
	ez.AssertAreEqual(dataDigest, recoveredHash)
}

func Test__GHashOverGCrypto__WithNonces(t *testing.T) {
	ez := ez.New(t)
	data := []byte("Hello, World!")
	key := []byte("This is a key")

	crypt := gcrypt.New(64)

	dataHash := ghash.NewWithParams(1, 64, 128, nil)
	dataHash.SetNonce([]byte("This is a nonce"))
	dataNonceState := dataHash.GetNonceState()
	encodedData := crypt.Encode(data)
	encodedDataBytes := crypt.EncodeToBytes(data)
	dataHash.AddBytes(encodedDataBytes)
	dataDigest := dataHash.GetDigest()

	keyHash := ghash.NewWithParams(1, 64, 128, nil)
	keyHash.SetNonce([]byte("This is a key nonce"))
	keyNonceState := keyHash.GetNonceState()

	encrypted := crypt.Encrypt(data, key)

	for i := 0; i < len(dataNonceState); i++ {
		dataNonceState[i].Add(keyNonceState[i])
	}

	encryptedHash := ghash.NewWithParams(1, 64, 128, nil)
	encryptedHash.SetNonceState(dataNonceState)
	encryptedHash.AddBytes(encrypted)
	encodedKey := crypt.ExpandKeyToBytes(key, len(encodedData))
	encryptedHash.RemoveBytes(encodedKey)

	encryptedHash.RemoveNonceState(keyNonceState)
	recoveredHash := encryptedHash.GetDigest()

	ez.AssertAreEqual(crypt.Decrypt(encrypted, key), data)
	ez.AssertAreEqual(dataDigest, recoveredHash)
}
