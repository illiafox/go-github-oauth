package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
)

// StoreState stores state, error can be only internal
func (m Memcached) StoreState(state string) error {
	return m.client.Set(&memcache.Item{
		Key:        state,
		Value:      nil,
		Flags:      0,
		Expiration: expire, // 1 hour
	})
}

// LookupState returns true if state existed, error can be only internal
func (m Memcached) LookupState(state string) (bool, error) {
	err := m.client.Delete(state)

	if err != nil {
		if err == memcache.ErrCacheMiss { // Not Found

			return false, nil
		}

		return false, err
	}

	return true, nil
}
