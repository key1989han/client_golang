package prometheus

import (
    "testing"
    "time"
)

func TestEvictableHistogramVec_Evict(t *testing.T) {
    vec := &EvictableHistogramVec{
        entries: make(map[string]*evictEntry),
        ttl:     100 * time.Millisecond,
    }

    vec.Get("test1")
    vec.Get("test2")

    time.Sleep(150 * time.Millisecond)

    evicted := vec.Evict()
    if evicted != 2 {
        t.Errorf("Expected 2 evicted, got %d", evicted)
    }
}
