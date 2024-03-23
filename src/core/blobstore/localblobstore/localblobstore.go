package localblobstore

import (
	"drm-blockchain/src/core/blobstore"
	"drm-blockchain/src/di"
	errorutils "drm-blockchain/src/utils/error"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

type LocalBlobStore struct {
	basePath string
}

func DIFactory(*di.DIContext) blobstore.BlobStore {
	return LocalBlobStore{basePath: "/tmp/drm-blockchain"}
}

func (bs LocalBlobStore) Exists(relPath string) (bool, error) {
	absPath, err := bs.prepare(relPath)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(absPath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func (bs LocalBlobStore) OpenRead(relPath string) (io.ReadCloser, error) {
	absPath, err := bs.prepare(relPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (bs LocalBlobStore) OpenWrite(relPath string) (io.WriteCloser, error) {
	absPath, err := bs.prepare(relPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(absPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (bs LocalBlobStore) Get(relPath string) ([]byte, error) {
	rd, err := bs.OpenRead(relPath)

	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(rd)

	if err != nil {
		return nil, err
	}

	err = rd.Close()
	return bytes, err
}

func (bs LocalBlobStore) Write(relPath string, content []byte) error {
	wr, err := bs.OpenWrite(relPath)

	if err != nil {
		return err
	}

	n, err := wr.Write(content)

	if err != nil {
		return err
	}

	if n != len(content) {
		return errorutils.Newf("Expected to write %d bytes, but could only write %d", len(content), n)
	}

	return wr.Close()
}

func (bs LocalBlobStore) prepare(relPath string) (string, error) {
	absPath := path.Join(bs.basePath, filepath.FromSlash(relPath))
	dirPath, _ := path.Split(absPath)
	return absPath, os.MkdirAll(dirPath, os.ModePerm)
}
