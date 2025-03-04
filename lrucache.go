package lrucache

import (
	"sync"
)

type LRUCache[K comparable, T any] interface {
	Get(key K) T
	Put(key K, value T)
	Delete(key K)
	Len() int
	Cap() int
	Keys() []K
}

type LRUCacheEntry[K comparable, T any] struct {
	key   K
	value T
}

type LRUCacheImpl[K comparable, T any] struct {
	len   int
	cap   int
	cache map[K]*Node[LRUCacheEntry[K, T]]
	list  DoubleLinkedList[LRUCacheEntry[K, T]]
	mu    sync.Mutex
}

func NewLRUCache[K comparable, T any](cap int) *LRUCacheImpl[K, T] {
	return &LRUCacheImpl[K, T]{
		cache: make(map[K]*Node[LRUCacheEntry[K, T]]),
		list:  NewDoubleLinkedList[LRUCacheEntry[K, T]](cap),
		cap:   cap,
		mu:    sync.Mutex{},
	}
}

func (lru *LRUCacheImpl[K, T]) Get(key K) (T, bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if node, ok := lru.cache[key]; ok {
		lru.list.MoveToHead(node)
		return node.data.value, true
	}

	var zero T
	return zero, false
}

func (lru *LRUCacheImpl[K, T]) Put(key K, value T) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if node, ok := lru.cache[key]; ok {
		node.data = LRUCacheEntry[K, T]{key: key, value: value}
		lru.list.MoveToHead(node)
		return
	}

	if lru.cap == lru.len {
		delete(lru.cache, lru.list.Tail().data.key)
		lru.list.RemoveTail()
		lru.len--
	}
	lru.list.Push(LRUCacheEntry[K, T]{key: key, value: value})
	lru.cache[key] = lru.list.Head()
	lru.len++
}

func (lru *LRUCacheImpl[K, T]) Delete(key K) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if node, ok := lru.cache[key]; ok {
		lru.list.Remove(node)
		delete(lru.cache, key)
	}
}

func (lru *LRUCacheImpl[T, K]) String() string {
	return lru.list.String()
}

func (lru *LRUCacheImpl[K, T]) Len() int {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.len
}

func (lru *LRUCacheImpl[K, T]) Cap() int {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.cap
}

func (lru *LRUCacheImpl[K, T]) Keys() []K {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	keys := make([]K, 0, lru.len)
	for node := lru.list.Head(); node.next != nil; node = node.next {
		keys = append(keys, node.data.key)
	}
	return keys
}
