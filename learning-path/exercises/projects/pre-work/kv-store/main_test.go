package main

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestStoreBasicOperations(t *testing.T) {
	store := NewStore("")

	// Test SET and GET
	store.Set("name", "Alice")
	if val, ok := store.Get("name"); !ok || val != "Alice" {
		t.Errorf("Get() = %q, %v; want %q, true", val, ok, "Alice")
	}

	// Test non-existent key
	if _, ok := store.Get("nonexistent"); ok {
		t.Error("Get() for nonexistent key should return false")
	}

	// Test DELETE
	if !store.Delete("name") {
		t.Error("Delete() should return true for existing key")
	}
	if _, ok := store.Get("name"); ok {
		t.Error("Get() after Delete() should return false")
	}

	// Test DELETE on non-existent key
	if store.Delete("nonexistent") {
		t.Error("Delete() should return false for nonexistent key")
	}
}

func TestStoreExists(t *testing.T) {
	store := NewStore("")

	store.Set("key1", "value1")

	if !store.Exists("key1") {
		t.Error("Exists() should return true for existing key")
	}

	if store.Exists("key2") {
		t.Error("Exists() should return false for nonexistent key")
	}
}

func TestStoreKeys(t *testing.T) {
	store := NewStore("")

	store.Set("user:1", "Alice")
	store.Set("user:2", "Bob")
	store.Set("admin:1", "Charlie")

	tests := []struct {
		pattern string
		want    int
	}{
		{"*", 3},
		{"user:*", 2},
		{"admin:*", 1},
		{"none:*", 0},
	}

	for _, tt := range tests {
		keys := store.Keys(tt.pattern)
		if len(keys) != tt.want {
			t.Errorf("Keys(%q) returned %d keys, want %d", tt.pattern, len(keys), tt.want)
		}
	}
}

func TestStoreSize(t *testing.T) {
	store := NewStore("")

	if store.Size() != 0 {
		t.Errorf("Size() = %d, want 0", store.Size())
	}

	store.Set("key1", "value1")
	store.Set("key2", "value2")

	if store.Size() != 2 {
		t.Errorf("Size() = %d, want 2", store.Size())
	}

	store.Delete("key1")

	if store.Size() != 1 {
		t.Errorf("Size() = %d, want 1", store.Size())
	}
}

func TestStoreClear(t *testing.T) {
	store := NewStore("")

	store.Set("key1", "value1")
	store.Set("key2", "value2")

	store.Clear()

	if store.Size() != 0 {
		t.Errorf("Size() after Clear() = %d, want 0", store.Size())
	}
}

func TestStoreSnapshotAndLoad(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "test.json")

	// Create store and add data
	store1 := NewStore(tmpfile)
	store1.Set("key1", "value1")
	store1.Set("key2", "value2")

	// Save snapshot
	if err := store1.Snapshot(); err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}

	// Create new store and load
	store2 := NewStore(tmpfile)
	if err := store2.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify data
	if val, ok := store2.Get("key1"); !ok || val != "value1" {
		t.Errorf("After Load(), Get(key1) = %q, %v; want %q, true", val, ok, "value1")
	}

	if val, ok := store2.Get("key2"); !ok || val != "value2" {
		t.Errorf("After Load(), Get(key2) = %q, %v; want %q, true", val, ok, "value2")
	}

	if store2.Size() != 2 {
		t.Errorf("After Load(), Size() = %d, want 2", store2.Size())
	}
}

func TestStoreConcurrentReads(t *testing.T) {
	store := NewStore("")
	store.Set("key1", "value1")

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if val, ok := store.Get("key1"); !ok || val != "value1" {
					t.Errorf("Concurrent Get() = %q, %v; want %q, true", val, ok, "value1")
				}
			}
		}()
	}

	wg.Wait()
}

func TestStoreConcurrentWrites(t *testing.T) {
	store := NewStore("")

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := "key"
			value := "value"
			for j := 0; j < 100; j++ {
				store.Set(key, value)
				store.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// Should complete without race conditions
}

func TestStoreConcurrentMixed(t *testing.T) {
	store := NewStore("")

	var wg sync.WaitGroup
	numReaders := 50
	numWriters := 50

	// Readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				store.Keys("*")
				store.Size()
				store.Get("key1")
			}
		}()
	}

	// Writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				store.Set("key1", "value1")
				store.Delete("key1")
			}
		}(i)
	}

	wg.Wait()
}

// Benchmarks

func BenchmarkStoreGet(b *testing.B) {
	store := NewStore("")
	store.Set("key", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get("key")
	}
}

func BenchmarkStoreSet(b *testing.B) {
	store := NewStore("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Set("key", "value")
	}
}

func BenchmarkStoreDelete(b *testing.B) {
	store := NewStore("")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Set("key", "value") // Setup
		store.Delete("key")
	}
}

func BenchmarkStoreConcurrentReads(b *testing.B) {
	store := NewStore("")
	store.Set("key", "value")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Get("key")
		}
	})
}

func BenchmarkStoreConcurrentWrites(b *testing.B) {
	store := NewStore("")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set("key", "value")
		}
	})
}

func BenchmarkStoreKeys(b *testing.B) {
	store := NewStore("")
	for i := 0; i < 1000; i++ {
		store.Set("key", "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Keys("*")
	}
}

func BenchmarkStoreSnapshot(b *testing.B) {
	tmpfile := filepath.Join(b.TempDir(), "bench.json")
	store := NewStore(tmpfile)

	// Add 10000 entries
	for i := 0; i < 10000; i++ {
		store.Set("key", "value")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := store.Snapshot(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStoreLoad(b *testing.B) {
	tmpfile := filepath.Join(b.TempDir(), "bench.json")
	store := NewStore(tmpfile)

	// Add 10000 entries and save
	for i := 0; i < 10000; i++ {
		store.Set("key", "value")
	}
	if err := store.Snapshot(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newStore := NewStore(tmpfile)
		if err := newStore.Load(); err != nil {
			b.Fatal(err)
		}
	}
}

// TODO: Add more tests
// - Test error conditions (invalid file, corrupted data)
// - Test large datasets
// - Test memory usage
// - Test with different data patterns
