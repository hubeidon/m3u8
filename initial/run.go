package initial

import "sync/atomic"

func Run() {

	go CompositeVideo()

}

// 使用atomic解决map并发写
func mapWrite(m *map[string]*int64, key string, delta int64) {
	// m["1"] = 10
	if (*m)[key] == nil {
		(*m)[key] = new(int64)
	}
	atomic.AddInt64((*m)[key], delta)
}
