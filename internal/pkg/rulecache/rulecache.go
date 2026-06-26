package rulecache

import "sync"

var Cache = struct {
	sync.RWMutex
	Group map[string]string
}{
	Group: make(map[string]string),
}

func SetGroups(groups map[string]string) {
	Cache.Lock()
	Cache.Group = groups
	Cache.Unlock()
}

func GetGroupChannels(name string) (string, bool) {
	Cache.RLock()
	defer Cache.RUnlock()
	ch, ok := Cache.Group[name]
	return ch, ok
}
