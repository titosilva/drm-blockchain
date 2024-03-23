package blobstore

import "io"

type BlobStore interface {
	Exists(path string) (bool, error)

	OpenRead(path string) (io.ReadCloser, error)
	OpenWrite(path string) (io.WriteCloser, error)

	Get(path string) ([]byte, error)
	Write(path string, content []byte) error
}
