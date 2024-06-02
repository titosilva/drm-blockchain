package lthash_test

import (
	"crypto/rand"
	"drm-blockchain/src/crypto/hash/lthash"
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

func Benchmark__LtHash__1GB(b *testing.B) {
	lt := lthash.NewDirect(500, 128, 1<<12, nil)
	bs, err := generateRandomBytes(1 << 30)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lt.Reset()

		if err != nil {
			b.Error(err)
		}

		lt.ComputeDigest(bs)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/hash")
}

func Benchmark__LtHash__1MB(b *testing.B) {
	lt := lthash.NewDirect(512, 128, 256, nil)
	bs, err := generateRandomBytes(1 << 20)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lt.Reset()

		if err != nil {
			b.Error(err)
		}

		lt.ComputeDigest(bs)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/hash")
}

func Benchmark__LtHash__1kB(b *testing.B) {
	lt := lthash.NewDirect(512, 128, 256, nil)
	bs, err := generateRandomBytes(1 << 10)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lt.Reset()

		if err != nil {
			b.Error(err)
		}

		lt.ComputeDigest(bs)
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/hash")
}
