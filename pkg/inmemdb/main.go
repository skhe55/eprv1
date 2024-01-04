package inmemdb

import "sync"

type InmemDB struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewInmemDB() *InmemDB {
	return &InmemDB{
		data: make(map[string]interface{}),
	}
}

func (d *InmemDB) Set(key string, value interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[key] = value
}

func (d *InmemDB) Get(key string) (interface{}, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	value, ok := d.data[key]
	return value, ok
}

func (d *InmemDB) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.data, key)
}
