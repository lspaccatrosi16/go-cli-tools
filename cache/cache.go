package cache

import (
	"time"

	"github.com/lspaccatrosi16/go-cli-tools/internal/pkgError"
	"github.com/lspaccatrosi16/go-libs/gbin"
)

var wrap = pkgError.WrapErrorFactory("cache")

type CacheItem struct {
	Data    []byte
	Created int64
}

func (c *CacheItem) IsValid(sec int64) bool {
	now := time.Now().Unix()
	return now-c.Created <= sec
}

func (c *CacheItem) Encode() ([]byte, error) {
	encoder := gbin.NewEncoder[CacheItem]()
	d, err := encoder.Encode(c)
	if err != nil {
		return nil, wrap(err)
	}
	return d, nil
}

func (c *CacheItem) Decode(b []byte) error {
	decoder := gbin.NewDecoder[CacheItem]()
	d, err := decoder.Decode(b)
	if err != nil {
		return wrap(err)
	}
	*c = *d
	return nil
}

func CreateCacheItem(data []byte) *CacheItem {
	now := time.Now().Unix()
	return &CacheItem{
		Data:    data,
		Created: now,
	}
}
