package hw04lrucache

type Key string

type CacheItem struct {
	key   Key
	Value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (cache lruCache) Get(key Key) (interface{}, bool) {
	if cache.items[key] != nil {
		listItem := cache.items[key]
		cache.queue.MoveToFront(listItem)
		return listItem.Value, true
	}

	return nil, false
}

func (cache lruCache) Set(key Key, value interface{}) bool {
	isHit := false
	if cache.items[key] == nil {
		if cache.queue.Len() >= cache.capacity {
			t := cache.queue.Back()
			cache.queue.Remove(t)
			deletedKey := t.Value.(CacheItem).key
			delete(cache.items, deletedKey)
		}

		cache.queue.PushFront(CacheItem{key, value})
		cache.items[key] = cache.queue.Front()
	} else {
		listItem := cache.items[key]
		listItem.Value = CacheItem{key, value}
		cache.queue.MoveToFront(listItem)
		isHit = true
	}

	return isHit
}

func (cache lruCache) Clear() {
	for key, value := range cache.items {
		cache.queue.Remove(value)
		delete(cache.items, key)
	}
}

func NewCache(capacity int) Cache {
	return lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
