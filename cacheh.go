package cacheh

import (
	"net/url"
	"strings"
)

type Cache interface {
	Get(key string) ([]byte, error) // returns nil if key not found
	Set(key string, value []byte) error
	Delete(key string) error // not an error if the key was not found

	WithKeyPrefix(keyPrefix string) Cache // get a Cache scoped to a key prefix
}

func GetDirCacheDsn(dir string) string {
	return "dir:" + dir
}

func GetDirCacheWithGzipDsn(dir string) string {
	return GetDirCacheDsn(dir) + "?gzip=1"
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
		return newFileCache(parsedDsn.rest, parsedDsn.query)
	default:
		return nil, ErrCacheInit{"unknown DSN kind: " + parsedDsn.kind, nil}
	}
}

const dsnSep = ":"

type parsedDsn struct {
	kind  string
	rest  string
	query map[string]string
}

func getParsedDsn(dsn string) (*parsedDsn, error) {
	parts := strings.SplitN(dsn, dsnSep, 2)

	if len(parts) != 2 {
		return nil, ErrCacheInit{"invalid DSN: " + dsn, nil}
	}

	pd := &parsedDsn{
		kind:  parts[0],
		rest:  parts[1],
		query: nil,
	}

	if strings.Contains(parts[1], "?") {
		spl := strings.SplitN(parts[1], "?", 2)
		pd.rest = spl[0]
		query, err := url.ParseQuery(spl[1])
		if err != nil {
			return nil, ErrCacheInit{"invalid DSN query " + spl[1], err}
		}
		pd.query = make(map[string]string)
		for k, v := range query {
			pd.query[k] = v[0]
		}
	}

	return pd, nil
}
