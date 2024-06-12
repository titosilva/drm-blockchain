package gcrypt_test

import (
	"crypto/rand"
	"drm-blockchain/src/crypto/encryption/gcrypt"
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

func runEncryptBenchmark(b *testing.B, messageBitsize int, chunkSize uint) {
	g := gcrypt.New(uint64(chunkSize))
	bs, _ := generateRandomBytes(messageBitsize / 8)
	key, _ := generateRandomBytes(32)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.Encrypt(bs, key)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/encryption")
}

func runDecryptBenchmark(b *testing.B, messageBitSize int, chunkSize uint) {
	g := gcrypt.New(uint64(chunkSize))
	bs, _ := generateRandomBytes(messageBitSize / 8)
	key, _ := generateRandomBytes(32)
	encrypted := g.Encrypt(bs, key)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.Decrypt(encrypted, key)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/decryption")
}

func Benchmark__GCrypt__128__128(b *testing.B) {
	runEncryptBenchmark(b, 128, 128)
}

func Benchmark__GCrypt__192__128(b *testing.B) {
	runEncryptBenchmark(b, 192, 128)
}

func Benchmark__GCrypt__256__128(b *testing.B) {
	runEncryptBenchmark(b, 256, 128)
}

func Benchmark__GCrypt__Decrypt__128__128(b *testing.B) {
	runDecryptBenchmark(b, 128, 128)
}

func Benchmark__GCrypt__Decrypt__192__128(b *testing.B) {
	runDecryptBenchmark(b, 192, 128)
}

func Benchmark__GCrypt__Decrypt__256__128(b *testing.B) {
	runDecryptBenchmark(b, 256, 128)
}
