package file

import (
	"context"
	"sync"
	"time"

	"github.com/go-micro/plugins/v4/store/file"
	jsoniter "github.com/json-iterator/go"
	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/store"
)

// NewCache file cache
//
// path: cache file存储路径
func NewCache(path string) cache.Cache {
	items := make(map[string]cache.Item)
	s := file.NewStore(store.Database("cache"), store.Table("file"), file.DirOption(path))
	keys, _ := s.List()
	if keys != nil {
		for i := 0; i < len(keys); i++ {
			records, err := s.Read(keys[i])
			if err != nil {
				continue
			}
			if len(records) != 0 {
				item := cache.Item{}
				_ = jsoniter.Unmarshal(records[0].Value, &item)
				items[records[0].Key] = item
			}
		}
	}
	return &fileCache{
		items: items,
		store: s,
	}
}

type fileCache struct {
	opts cache.Options
	sync.RWMutex

	items map[string]cache.Item
	store store.Store
}

func (f *fileCache) Get(ctx context.Context, key string) (interface{}, time.Time, error) {
	f.RWMutex.RLock()
	defer f.RWMutex.RUnlock()

	item, found := f.items[key]
	if !found {
		return nil, time.Time{}, cache.ErrKeyNotFound
	}
	if item.Expired() {
		return nil, time.Time{}, cache.ErrItemExpired
	}

	return item.Value, time.Unix(0, item.Expiration), nil
}

func (f *fileCache) Put(ctx context.Context, key string, val interface{}, d time.Duration) error {
	var e int64
	if d == cache.DefaultExpiration {
		d = f.opts.Expiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	f.RWMutex.Lock()
	defer f.RWMutex.Unlock()

	item := cache.Item{
		Value:      val,
		Expiration: e,
	}
	buf, _ := jsoniter.Marshal(item)
	err := f.store.Write(&store.Record{
		Key:      key,
		Value:    buf,
		Metadata: nil,
		Expiry:   d,
	})
	if err != nil {
		return err
	}
	f.items[key] = item

	return nil
}

func (f *fileCache) Delete(ctx context.Context, key string) error {
	f.RWMutex.Lock()
	defer f.RWMutex.Unlock()

	_, found := f.items[key]
	if !found {
		return cache.ErrKeyNotFound
	}

	if err := f.store.Delete(key); err != nil {
		return err
	}
	delete(f.items, key)

	return nil
}

func (f *fileCache) String() string {
	return "file"
}
