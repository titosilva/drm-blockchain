package ghash_test

import (
	"drm-blockchain/src/crypto/hash/ghash"
	"testing"
)

func runBenchmark(b *testing.B, indexCount int, chunkCount uint, chunkSize uint) {
	g := ghash.NewWithParams(chunkCount, chunkSize, 64, nil)
	bs, err := generateRandomBytes(indexCount * int(chunkSize) / 8)
	b.ResetTimer()

	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		g.AddBytes(bs)
		g.GetDigest()
	}

	b.ReportMetric(float64(b.Elapsed().Milliseconds())/float64(b.N), "ms/hash")
}

func Benchmark__GHash__64__500__128(b *testing.B) {
	runBenchmark(b, 64, 500, 128)
}

func Benchmark__GHash__128__500__128(b *testing.B) {
	runBenchmark(b, 128, 500, 128)
}

func Benchmark__GHash__192__500__128(b *testing.B) {
	runBenchmark(b, 192, 500, 128)
}

func Benchmark__GHash__256__500__128(b *testing.B) {
	runBenchmark(b, 256, 500, 128)
}
