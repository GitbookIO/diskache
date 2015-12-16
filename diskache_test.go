package diskache

import (
	"bytes"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	TMP_DIR     = "tmp"
)

func randStringBytes(n int) (string, []byte) {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b), b
}

func cleanDir() {
	os.RemoveAll(TMP_DIR)
}

// Test Set-ting and Get-ting a value in cache
func TestSetGet(t *testing.T) {
	// Cleanup
	defer cleanDir()

	// Create an instance
	opts := &Opts{
		Directory: TMP_DIR,
	}

	dc, err := New(opts)
	if err != nil {
		t.Error("Expected to create a new instance of Diskache")
	}

	// Set a value in cache
	key, value := randStringBytes(16)
	err = dc.Set(key, value)
	if err != nil {
		t.Error("Expected to be able to set value in cache")
	}

	// Read from cache
	cached, inCache := dc.Get(key)
	if inCache {
		if comp := bytes.Compare(value, cached); comp != 0 {
			t.Error("Expected to get the same value that was set in cache")
		}
	} else {
		t.Error("Expected to get value from cache")
	}
}

// Test concurrent Set by multiple go routines
func TestConcurrent(t *testing.T) {
	// Cleanup
	defer cleanDir()

	// Create an instance
	opts := &Opts{
		Directory: TMP_DIR,
	}

	dc, err := New(opts)
	if err != nil {
		t.Error("Expected to create a new instance of Diskache")
	}

	// Set multiple times the same value
	key, value := randStringBytes(16)
	for i := 0; i < 1000; i++ {
		go func() {
			if err := dc.Set(key, value); err != nil {
				t.Error("Expected Diskache to handle concurrency")
			}
		}()
	}

	// Read from cache
	cached, inCache := dc.Get(key)
	if inCache {
		if comp := bytes.Compare(value, cached); comp != 0 {
			t.Error("Expected to get the same value that was set in cache")
		}
	} else {
		t.Error("Expected to get value from cache")
	}
}

// Benchmark Set operations
func BenchmarkSet(b *testing.B) {
	// Cleanup
	defer cleanDir()

	// Create an instance
	opts := &Opts{
		Directory: TMP_DIR,
	}

	dc, err := New(opts)
	if err != nil {
		b.Error("Expected to create a new instance of Diskache")
	}

	// Set multiple values in cache
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key, value := randStringBytes(16)
		err = dc.Set(key, value)
		if err != nil {
			b.Error("Expected to be able to set value in cache")
		}
	}
}

// Benchmark Get operations
func BenchmarkGet(b *testing.B) {
	// Cleanup
	defer cleanDir()

	// Create an instance
	opts := &Opts{
		Directory: TMP_DIR,
	}

	dc, err := New(opts)
	if err != nil {
		b.Error("Expected to create a new instance of Diskache")
	}

	// Set a value in cache
	key, value := randStringBytes(16)
	err = dc.Set(key, value)
	if err != nil {
		b.Error("Expected to be able to set value in cache")
	}

	// Read from cache
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dc.Get(key)
	}
}
