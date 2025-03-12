package pokecache 

import (
  "time"
  "sync"
)
type cacheEntry struct{
  createdAt time.Time
  val []byte
}

type Cache struct{
  cache map[string]cacheEntry
  m sync.Mutex
}

func NewCache(interval time.Duration) *Cache{
  cache := &Cache{
    cache: make(map[string]cacheEntry),
    m: sync.Mutex{},
  }
  
  go cache.readLoop(interval)

  return cache
}

func (c *Cache) Add(Key string, val []byte){
  c.m.Lock()
  c.cache[Key] = cacheEntry{
    createdAt:time.Now(),
    val: val,
  }
  c.m.Unlock()
}

func (c *Cache) Get(Key string)([]byte, bool){
  c.m.Lock()
  k,ok := c.cache[Key]
  c.m.Unlock()
  if ok{
    return k.val, true
  }
  return []byte{}, false
}

func (c *Cache) readLoop(interval time.Duration){
  ticker := time.NewTicker(interval)
  defer ticker.Stop()
  
  for{
   <- ticker.C //wait for next tick

  c.m.Lock()
  for k,v := range c.cache{
    elapsed := time.Since(v.createdAt)
    if elapsed > interval{
      delete(c.cache,k)
    }
  }
  c.m.Unlock()
}
}
