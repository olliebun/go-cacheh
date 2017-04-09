package cacheh

import "fmt"

type ErrCacheInit struct {
	msg string
	err error
}

func (e ErrCacheInit) Error() string {
	if e.err != nil {
		return fmt.Sprintf("ErrCacheInit: %s: %s", e.msg, e.err)
	} else {

		return fmt.Sprintf("ErrCacheInit: %s", e.msg)
	}
}

type ErrUnsafeCacheKey struct {
	key string
}

func (e ErrUnsafeCacheKey) Error() string {
	return fmt.Sprintf("ErrCacheUnsafeKey: %s is an unsafe key", e.key)
}

type ErrCacheOperation struct {
	op  string
	key string
	err error
}

func (e ErrCacheOperation) Error() string {
	return fmt.Sprintf("ErrCacheOperation: %s on %s failed: %s", e.op, e.key, e.err)
}
