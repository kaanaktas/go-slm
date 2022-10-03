package cache

import (
	"time"
)

type Cache interface {
	Get(k string) (any, bool)
	Set(k string, v any, d time.Duration)
	Flush()
}

const (
	NoExpiration time.Duration = -1
)
