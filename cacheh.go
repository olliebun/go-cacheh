package cacheh

import (
	"strings"
)

type Cache interface {
	Get(string) ([]byte, error) // returns nil if key not found
	Set(string, []byte) error
	Delete(string) error // not an error if the key was not found
}

func GetDirCacheDsn(dir string) string {
	return "dir:" + dir
}

// NewCache constructs a new cache based on the dsn.
//
// For example, file-based:
//
//  NewCache("dir:/home/$user/")
func NewCache(dsn string) (Cache, error) {
	parsedDsn, err := getParsedDsn(dsn)

	if err != nil {
		return nil, err
	}

	switch strings.ToLower(parsedDsn.kind) {
	case "dir":
		return newFileCache(parsedDsn.rest)
	default:
		return nil, ErrCacheInit{"unknown DSN kind: " + parsedDsn.kind, nil}
	}
}

const dsnSep = ":"

type parsedDsn struct {
	kind string
	rest string
}

func getParsedDsn(dsn string) (*parsedDsn, error) {
	parts := strings.SplitN(dsn, dsnSep, 2)

	if len(parts) != 2 {
		return nil, ErrCacheInit{"invalid DSN: " + dsn, nil}
	}

	return &parsedDsn{
		parts[0],
		parts[1],
	}, nil
}
