package ghash_test

import (
	"crypto/rand"
	"drm-blockchain/src/crypto/encryption/gcrypt"
	"drm-blockchain/src/crypto/hash/ghash"
	"drm-blockchain/src/math/uintp"
	"reflect"
	"testing"
)

func generateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)

	_, err := rand.Read(bytes)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func hashThenEncrypt(data []byte, key []byte, modulusBitsize uint64, block_size_bits int) ([]byte, []*uintp.UintP, []byte) {
	crypt := gcrypt.New(modulusBitsize)

	dataHash := ghash.NewWithParams(500, uint(modulusBitsize), block_size_bits/8, nil)
	dataHash.SetNonce([]byte("This is a nonce"))
	dataNonceState := dataHash.GetNonceState()
	encodedDataBytes := crypt.EncodeToBytes(data)
	dataHash.AddBytes(encodedDataBytes)

	keyHash := ghash.NewWithParams(500, uint(modulusBitsize), block_size_bits/8, nil)
	keyHash.SetNonce([]byte("This is a key nonce"))
	keyNonceState := keyHash.GetNonceState()

	encrypted := crypt.Encrypt(data, key)

	for i := 0; i < len(dataNonceState); i++ {
		dataNonceState[i].Add(keyNonceState[i])
	}

	return encrypted, dataNonceState, dataHash.GetDigest()
}

func hashEncrypted(encrypted []byte, encryptedNonceState []*uintp.UintP, modulusBitsize uint64, block_size_bits int) []*uintp.UintP {
	encryptedHash := ghash.NewWithParams(500, uint(modulusBitsize), block_size_bits/8, nil)
	encryptedHash.SetNonceState(encryptedNonceState)
	encryptedHash.AddBytes(encrypted)

	return encryptedHash.GetState()
}

func verifyHash(encryptedHashState []*uintp.UintP, encrypted []byte, key []byte, modulusBitsize uint64, block_size_bits int) {
	crypt := gcrypt.New(modulusBitsize)
	encodedKey := crypt.ExpandKeyToBytes(key, len(encrypted))

	encryptedHashObj := ghash.NewWithParams(500, uint(modulusBitsize), block_size_bits/8, nil)
	encryptedHashObj.SetNonceState(encryptedHashState)
	encryptedHashObj.RemoveBytes(encodedKey)
	encryptedHashObj.RemoveNonce([]byte("This is a key nonce"))
}

func decrypt(encrypted []byte, key []byte, modulusBitsize uint64) []byte {
	crypt := gcrypt.New(modulusBitsize)
	return crypt.Decrypt(encrypted, key)
}

func runPdpr(data []byte, key []byte, modulusBitsize uint64, block_size_bits int) {
	crypt := gcrypt.New(modulusBitsize)

	dataHash := ghash.NewWithParams(500, uint(modulusBitsize), block_size_bits/8, nil)
	dataHash.SetNonce([]byte("This is a nonce"))
	dataNonceState := dataHash.GetNonceState()
	encodedData := crypt.Encode(data)
	encodedDataBytes := crypt.EncodeToBytes(data)
	dataHash.AddBytes(encodedDataBytes)
	dataDigest := dataHash.GetDigest()

	keyHash := ghash.NewWithParams(500, uint(modulusBitsize), block_size_bits/8, nil)
	keyHash.SetNonce([]byte("This is a key nonce"))
	keyNonceState := keyHash.GetNonceState()

	encrypted := crypt.Encrypt(data, key)

	for i := 0; i < len(dataNonceState); i++ {
		dataNonceState[i].Add(keyNonceState[i])
	}

	encryptedHash := ghash.NewWithParams(500, uint(modulusBitsize), block_size_bits/8, nil)
	encryptedHash.SetNonceState(dataNonceState)
	encryptedHash.AddBytes(encrypted)
	encodedKey := crypt.ExpandKeyToBytes(key, len(encodedData))
	encryptedHash.RemoveBytes(encodedKey)

	encryptedHash.RemoveNonceState(keyNonceState)
	recoveredHash := encryptedHash.GetDigest()

	if !reflect.DeepEqual(dataDigest, recoveredHash) {
		panic("Data hash and recovered hash do not match")
	}
}

func Benchmark__FullPdpr__HelloWorld__128m__128b(b *testing.B) {
	data := []byte("Hello, World!")
	key := []byte("This is a key")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		runPdpr(data, key, 128, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/pdpr")
}

func Benchmark__FullPdpr__256bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		runPdpr(data, key, 128, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/pdpr")
}

func Benchmark__HashThenEncrypt__256bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashThenEncrypt(data, key, 128, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__HashEncrypted__256bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	encrypted, encryptedNonceState, _ := hashThenEncrypt(data, key, 128, 128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashEncrypted(encrypted, encryptedNonceState, 128, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__VerifyHash__256bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	encrypted, encryptedNonceState, _ := hashThenEncrypt(data, key, 128, 128)
	encryptedHashState := hashEncrypted(encrypted, encryptedNonceState, 128, 128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		verifyHash(encryptedHashState, encrypted, key, 128, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__Decrypt__256bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	encrypted, _, _ := hashThenEncrypt(data, key, 128, 128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decrypt(encrypted, key, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__HashThenEncrypt__128bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(16)
	key, _ := generateRandomBytes(16)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashThenEncrypt(data, key, 128, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__HashEncrypted__128bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(16)
	key, _ := generateRandomBytes(16)
	encrypted, encryptedNonceState, _ := hashThenEncrypt(data, key, 128, 128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashEncrypted(encrypted, encryptedNonceState, 128, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__VerifyHash__128bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(16)
	key, _ := generateRandomBytes(16)
	encrypted, encryptedNonceState, _ := hashThenEncrypt(data, key, 128, 128)
	encryptedHashState := hashEncrypted(encrypted, encryptedNonceState, 128, 128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		verifyHash(encryptedHashState, encrypted, key, 128, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__Decrypt__128bit__128m__128b(b *testing.B) {
	data, _ := generateRandomBytes(16)
	key, _ := generateRandomBytes(16)
	encrypted, _, _ := hashThenEncrypt(data, key, 128, 128)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decrypt(encrypted, key, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__HashThenEncrypt__256bit__128m__256b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashThenEncrypt(data, key, 128, 256)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__HashEncrypted__256bit__128m__256b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	encrypted, encryptedNonceState, _ := hashThenEncrypt(data, key, 128, 256)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashEncrypted(encrypted, encryptedNonceState, 128, 256)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__VerifyHash__256bit__128m__256b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	encrypted, encryptedNonceState, _ := hashThenEncrypt(data, key, 128, 256)
	encryptedHashState := hashEncrypted(encrypted, encryptedNonceState, 128, 256)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		verifyHash(encryptedHashState, encrypted, key, 128, 256)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__Decrypt__256bit__128m__256b(b *testing.B) {
	data, _ := generateRandomBytes(32)
	key, _ := generateRandomBytes(32)
	encrypted, _, _ := hashThenEncrypt(data, key, 128, 256)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decrypt(encrypted, key, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__HashThenEncrypt__128bit__128m__256b(b *testing.B) {
	data, _ := generateRandomBytes(16)
	key, _ := generateRandomBytes(16)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashThenEncrypt(data, key, 128, 256)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__HashEncrypted__128bit__128m__256b(b *testing.B) {
	data, _ := generateRandomBytes(16)
	key, _ := generateRandomBytes(16)
	encrypted, encryptedNonceState, _ := hashThenEncrypt(data, key, 128, 256)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hashEncrypted(encrypted, encryptedNonceState, 128, 256)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__VerifyHash__128bit__128m__256b(b *testing.B) {
	data, _ := generateRandomBytes(16)
	key, _ := generateRandomBytes(16)
	encrypted, encryptedNonceState, _ := hashThenEncrypt(data, key, 128, 256)
	encryptedHashState := hashEncrypted(encrypted, encryptedNonceState, 128, 256)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		verifyHash(encryptedHashState, encrypted, key, 128, 256)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}

func Benchmark__Decrypt__128bit__128m__256b(b *testing.B) {
	data, _ := generateRandomBytes(16)
	key, _ := generateRandomBytes(16)
	encrypted, _, _ := hashThenEncrypt(data, key, 128, 256)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		decrypt(encrypted, key, 128)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms")
}
