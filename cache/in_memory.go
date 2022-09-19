package cache

import (
	gocache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type inMemory struct {
	cache *gocache.Cache
}

func (i *inMemory) Get(k string) (interface{}, bool) {
	return i.cache.Get(k)
}

func (i *inMemory) Set(k string, v interface{}, d time.Duration) {
	i.cache.Set(k, v, d*time.Minute)
}

func (i *inMemory) Flush() {
	i.cache.Flush()
}

var onceInMem sync.Once

func (i *inMemory) initiateInMemory() {
	onceInMem.Do(func() {
		i.cache = gocache.New(5*time.Minute, 10*time.Minute)
	})
}

var inMem inMemory

func NewInMemory() Cache {
	inMem.initiateInMemory()
	return &inMem
}
