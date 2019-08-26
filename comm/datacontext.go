package comm

import (
	"sync"
	"time"
)

//保存键值数据
type DataContext struct {
	mutex           sync.RWMutex
	DataStore       map[interface{}]interface{}
	DataStoreExpire map[interface{}]int64
}

// Set stores a value for a given key in a given request.
func (d *DataContext) Set(key, val interface{}) {
	d.mutex.Lock()
	if d.DataStore == nil {
		d.DataStore = make(map[interface{}]interface{})
		d.DataStoreExpire = make(map[interface{}]int64)
	}
	d.DataStore[key] = val
	d.DataStoreExpire[key] = time.Now().Unix()
	d.mutex.Unlock()
}

//Refresh refresh expire time
func (d *DataContext) Refresh(key interface{}) {
	d.mutex.Lock()
	d.DataStoreExpire[key] = time.Now().Unix()
	d.mutex.Unlock()
}

// Get returns a value stored for a given key in a given request.
func (d *DataContext) Get(key interface{}) interface{} {
	d.mutex.RLock()
	if ctx := d.DataStore[key]; ctx != nil {
		d.mutex.RUnlock()
		return ctx
	}
	d.mutex.RUnlock()
	return nil
}

// GetOk returns stored value and presence state like multi-value return of map access.
func (d *DataContext) GetOk(key interface{}) (interface{}, bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if value, ok := d.DataStore[key]; ok {
		return value, ok
	}
	return nil, false
}

// Delete removes a value stored for a given key
func (d *DataContext) Delete(key interface{}) {
	d.mutex.Lock()
	if d.DataStore[key] != nil {
		delete(d.DataStore, key)
		delete(d.DataStoreExpire, key)
	}
	d.mutex.Unlock()
}

// Purge removes request data stored for longer than maxAge, in seconds.
// It returns the amount of requests removed.
//
// If maxAge <= 0, all request data is removed.
//清理
func (d *DataContext) Purge(key interface{}, maxAge int) int {
	return d.PurgeWithFunc(key, maxAge, nil)
}

func (d *DataContext) PurgeWithFunc(key interface{}, maxAge int, OnDelete func(k interface{}, v interface{})) int {
	count := 0
	if maxAge <= 0 {
		count = len(d.DataStore)
		d.DataStore = make(map[interface{}]interface{})
		d.DataStoreExpire = make(map[interface{}]int64)
	} else {
		min := time.Now().Unix() - int64(maxAge)
		for key, _ := range d.DataStore {
			if d.DataStoreExpire[key] < min {
				if OnDelete != nil {
					OnDelete(key, d.DataStore[key])
				}
				d.Delete(key)
				count++
			}
		}
	}
	//d.mutex.Unlock()
	return count
}
