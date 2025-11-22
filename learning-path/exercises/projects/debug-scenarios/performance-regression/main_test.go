package perfregression

import (
	"testing"
)

func BenchmarkExecuteSlow(b *testing.B) {
	data := make([]Row, 1000)
	for i := range data {
		data[i] = Row{"id": i, "value": "test"}
	}

	executor := &Executor{data: data}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = executor.Execute()
	}
}

func BenchmarkCacheSlow(b *testing.B) {
	cache := &Cache{
		data: make(map[string]interface{}),
	}

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cache.Put(string(rune(i)), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(string(rune(i % 1000)))
			i++
		}
	})
}

func BenchmarkJoinSlow(b *testing.B) {
	left := make([]Row, 100)
	right := make([]Row, 100)

	for i := range left {
		left[i] = Row{"id": i}
		right[i] = Row{"id": i, "value": i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashJoin(left, right)
	}
}

// After fixes, benchmark should show significant improvement
// Use: benchstat old.txt new.txt
