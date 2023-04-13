package utils

type SafeMap[K, V any] struct {
	channel chan *KVEntry[K, V]
	core    map[any]V
}

func NewSafeMap[K any, V any]() *SafeMap[K, V] {
	m := &SafeMap[K, V]{
		channel: make(chan *KVEntry[K, V]),
		core:    make(map[any]V),
	}
	go func() {
		for {
			select {
			case c := <-m.channel:
				m.core[c.key] = c.value
				break
			}
		}
	}()
	return m
}

type KVEntry[K, V any] struct {
	key   K
	value V
}

func (m *SafeMap[K, V]) Put(k K, v V) {
	item := &KVEntry[K, V]{key: k, value: v}
	m.channel <- item
}

func (m *SafeMap[K, V]) Get(k K) V {
	return m.core[k]
}
