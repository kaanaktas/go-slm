package cache

import (
	"time"
)

type Cache interface {
	Get(k string) (interface{}, bool)
	Set(k string, v interface{}, d time.Duration) error
}

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)
