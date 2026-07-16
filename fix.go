package prometheus

import "sync"
import "time"

// EvictableHistogramVec adds TTL-based eviction to HistogramVec
type EvictableHistogramVec struct {
    mu       sync.RWMutex
    entries  map[string]*evictEntry
    ttl      time.Duration
}

type evictEntry struct {
    lastAccess time.Time
    histogram  *Histogram
}

// Get retrieves or creates a histogram, resetting TTL
func (e *EvictableHistogramVec) Get(labels string) *Histogram {
    e.mu.Lock()
    defer e.mu.Unlock()

    if entry, ok := e.entries[labels]; ok {
        entry.lastAccess = time.Now()
        return entry.histogram
    }

    e.entries[labels] = &evictEntry{
        lastAccess: time.Now(),
    }
    return e.entries[labels].histogram
}

// Evict removes entries older than TTL
func (e *EvictableHistogramVec) Evict() int {
    e.mu.Lock()
    defer e.mu.Unlock()

    now := time.Now()
    evicted := 0
    for k, v := range e.entries {
        if now.Sub(v.lastAccess) > e.ttl {
            delete(e.entries, k)
            evicted++
        }
    }
    return evicted
}
