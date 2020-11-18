package xor

import "sync"

// 与或门 配置

var (
	cache = make(map[string]bool)
	mutex sync.RWMutex
)

//X 获取
func X(key string) bool {
	mutex.RLock()
	_, ok := cache[key]
	mutex.RUnlock()
	return ok
}

//O 设置存在为true
func O(key string) {
	mutex.Lock()
	cache[key] = true
	mutex.Unlock()
}

//R 删除值
func R(key string) {
	mutex.Lock()
	delete(cache, key)
	mutex.Unlock()
}
