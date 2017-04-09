package cacheh

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/kennygrant/sanitize"
)

func initFileCacheDir(dir string) error {
	fi, err := os.Stat(dir)

	if err != nil && !os.IsNotExist(err) {
		return ErrCacheInit{"could not stat " + dir, err}
	} else if err != nil {
		// create the dir
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return ErrCacheInit{"could not create " + dir, err}
		}
	} else if err == nil && !fi.IsDir() {
		return ErrCacheInit{fmt.Sprintf("%s exists but is not a directory", dir), nil}
	}

	return nil
}

func newFileCache(dir string) (Cache, error) {
	err := initFileCacheDir(dir)
	if err != nil {
		return nil, err
	}

	return &fileCache{dir, new(sync.RWMutex)}, nil
}

type fileCache struct {
	dir string
	*sync.RWMutex
}

func (fc *fileCache) keyPath(key string) string {
	return filepath.Join(fc.dir, key)
}

func (fc *fileCache) Get(key string) ([]byte, error) {
	// it is the caller's responsibility to provide sanitized strings - but
	// we still don't want to allow unsafed strings, so we will make sure
	// it's already sanitized and return an error if not.
	if !isSafeFileCacheKey(key) {
		return nil, ErrUnsafeCacheKey{key}
	}

	fc.RLock()
	defer fc.RUnlock()

	path := fc.keyPath(key)

	f, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, ErrCacheOperation{"read", key, err}
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, ErrCacheOperation{"read", key, err}
	}
	return b, nil
}

func (fc *fileCache) Set(key string, val []byte) error {
	// it is the caller's responsibility to provide sanitized strings - but
	// we still don't want to allow unsafed strings, so we will make sure
	// it's already sanitized and return an error if not.
	if !isSafeFileCacheKey(key) {
		return ErrUnsafeCacheKey{key}
	}
	fc.Lock()
	defer fc.Unlock()

	path := fc.keyPath(key)

	err := ioutil.WriteFile(path, val, 0600)

	if err != nil {
		return ErrCacheOperation{"write", key, err}
	}

	return nil
}

func (fc *fileCache) Delete(key string) error {
	// it is the caller's responsibility to provide sanitized strings - but
	// we still don't want to allow unsafed strings, so we will make sure
	// it's already sanitized and return an error if not.
	if !isSafeFileCacheKey(key) {
		return ErrUnsafeCacheKey{key}
	}

	fc.Lock()
	defer fc.Unlock()

	path := fc.keyPath(key)

	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return ErrCacheOperation{"delete", key, err}
	}
	return err
}

func isSafeFileCacheKey(key string) bool {
	ext := filepath.Ext(key)
	if ext != "" {
		key = key[:len(key)-len(ext)]
	}

	return sanitizeFilename(key) == key
}

// sanitizeFilename sanitises a string so that it can be safely used as a
// filename basename.
func sanitizeFilename(s string) string {
	return sanitize.BaseName(s)
}
