package cachedregexp

import (
	"regexp"
	"sync"
)

//nolint:gochecknoglobals // this is the whole point of being a cache
var cache sync.Map

func MustCompile(exp string) *regexp.Regexp {
	compiled, ok := cache.Load(exp)
	if !ok {
		compiled, _ = cache.LoadOrStore(exp, regexp.MustCompile(exp))
	}

	return compiled.(*regexp.Regexp)
}
