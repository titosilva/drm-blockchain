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

func runBenchmark(b *testing.B, sizeOfFile int, blockSize int, chunkCount uint, chunkSize uint) {
	lt := lthash.NewDirect(chunkCount, chunkSize, blockSize, nil)
	bs, err := generateRandomBytes(sizeOfFile)
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

func Benchmark__LtHash__1GB__4kB__500__128(b *testing.B) {
	runBenchmark(b, 1<<30, 1<<12, 500, 128)
}

func Benchmark__LtHash__1GB__4kB__500__256(b *testing.B) {
	runBenchmark(b, 1<<30, 1<<12, 500, 256)
}

func Benchmark__LtHash__1GB__1kB__500__128(b *testing.B) {
	runBenchmark(b, 1<<30, 1<<10, 500, 128)
}

func Benchmark__LtHash__1GB__1kB__500__256(b *testing.B) {
	runBenchmark(b, 1<<30, 1<<10, 500, 256)
}

func Benchmark__LtHash__1MB__4kB__500__128(b *testing.B) {
	runBenchmark(b, 1<<20, 1<<12, 500, 128)
}

func Benchmark__LtHash__1MB__4kB__500__256(b *testing.B) {
	runBenchmark(b, 1<<20, 1<<12, 500, 256)
}

func Benchmark__LtHash__1MB__1kB__500__128(b *testing.B) {
	runBenchmark(b, 1<<20, 1<<10, 500, 128)
}

func Benchmark__LtHash__1MB__1kB__500__256(b *testing.B) {
	runBenchmark(b, 1<<20, 1<<10, 500, 256)
}

func Benchmark__LtHash__1kB__512__500__128(b *testing.B) {
	runBenchmark(b, 1<<10, 512, 500, 128)
}

func Benchmark__LtHash__1kB__512__500__256(b *testing.B) {
	runBenchmark(b, 1<<10, 512, 500, 256)
}
