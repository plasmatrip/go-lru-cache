package golrucache

import (
	"sync"
)

type DoubleLinkedList[T any] interface {
	Head() *node[T]
	Tail() *node[T]
	Push(value T)
	Remove(node *node[T])
	removeTail()
	moveToHead(node *node[T])
	String() string
}

type LRUCache[K comparable, T any] interface {
	Get(key K) T
	Put(key K, value T)
	Delete(key K)
	Len() int
	Cap() int
	Keys() []K
}

type lruCacheEntry[K comparable, T any] struct {
	key   K
	value T
}

type lruCache[K comparable, T any] struct {
	len   int
	cap   int
	cache map[K]*node[lruCacheEntry[K, T]]
	list  DoubleLinkedList[lruCacheEntry[K, T]]
	mu    sync.Mutex
}

// NewLRUCache creates a new lruCache with a given capacity. The capacity is the
// maximum number of items that the cache can store. The cache will evict the
// least recently used item when it reaches capacity.
func NewLRUCache[K comparable, T any](cap int) LRUCache[K, T] {
	return &lruCache[K, T]{
		cache: make(map[K]*node[lruCacheEntry[K, T]]),
		list:  NewDoublyLinkedList[lruCacheEntry[K, T]](cap),
		cap:   cap,
		mu:    sync.Mutex{},
	}
}

// Get retrieves the value associated with the given key from the cache. If the
// key is found, the associated value is returned along with true. If the key is
// not found, the zero value for the type T is returned along with false.
func (lru *lruCache[K, T]) Get(key K) T {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if node, ok := lru.cache[key]; ok {
		lru.list.moveToHead(node)
		return node.data.value
	}

	var zero T
	return zero
}

func (lru *lruCache[K, T]) Put(key K, value T) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if node, ok := lru.cache[key]; ok {
		node.data = lruCacheEntry[K, T]{key: key, value: value}
		lru.list.moveToHead(node)
		return
	}

	if lru.cap == lru.len {
		delete(lru.cache, lru.list.Tail().data.key)
		lru.list.removeTail()
		lru.len--
	}
	lru.list.Push(lruCacheEntry[K, T]{key: key, value: value})
	lru.cache[key] = lru.list.Head()
	lru.len++
}

func (lru *lruCache[K, T]) Delete(key K) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if node, ok := lru.cache[key]; ok {
		lru.list.Remove(node)
		delete(lru.cache, key)
	}
}

func (lru *lruCache[T, K]) String() string {
	return lru.list.String()
}

func (lru *lruCache[K, T]) Len() int {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.len
}

func (lru *lruCache[K, T]) Cap() int {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.cap
}

func (lru *lruCache[K, T]) Keys() []K {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	keys := make([]K, 0, lru.len)
	for node := lru.list.Head(); node.next != nil; node = node.next {
		keys = append(keys, node.data.key)
	}
	return keys
}
