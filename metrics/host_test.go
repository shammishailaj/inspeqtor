package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectHostMetrics(t *testing.T) {
	store := NewHostStore()
	err := CollectHostMetrics(store, "proc")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, store.Get("cpu", ""), 0)
	assert.Equal(t, store.Get("cpu", "user"), 0)
	assert.Equal(t, store.Get("cpu", "system"), 0)
	assert.Equal(t, store.Get("cpu", "iowait"), 0)
	assert.Equal(t, store.Get("cpu", "steal"), 0)
	assert.Equal(t, store.Get("load", "1"), 2)
	assert.Equal(t, store.Get("load", "5"), 3)
	assert.Equal(t, store.Get("load", "15"), 5)
	assert.Equal(t, store.Get("swap", ""), 2)

	err = CollectHostMetrics(store, "proc2")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, store.Get("cpu", ""), 10)
	assert.Equal(t, store.Get("cpu", "user"), 1)
	assert.Equal(t, store.Get("cpu", "system"), 2)
	assert.Equal(t, store.Get("cpu", "iowait"), 3)
	assert.Equal(t, store.Get("cpu", "steal"), 4)
	assert.Equal(t, store.Get("load", "1"), 2)
	assert.Equal(t, store.Get("load", "5"), 3)
	assert.Equal(t, store.Get("load", "15"), 5)
	assert.Equal(t, store.Get("swap", ""), 2)
}

func TestCollectRealHostMetrics(t *testing.T) {
	store := NewHostStore()
	err := CollectHostMetrics(store, "/proc")
	if err != nil {
		t.Fatal(err)
	}
	// Can't really know what we'll collect so we'll check for non-zero.
	assert.True(t, store.Get("load", "1") > 0)
	assert.True(t, store.Get("load", "5") > 0)
	assert.True(t, store.Get("load", "15") > 0)
	assert.True(t, store.Get("swap", "") > 0)
}

func TestCollectDiskMetrics(t *testing.T) {
	store := NewHostStore()
	err := collectDisk("fixtures/df.linux.txt", store)
	if err != nil {
		t.Error(err)
	}
	if store.Get("disk", "/") != 17 {
		t.Error("Unexpected results: %v", store.Get("disk", "/"))
	}
	if store.Get("disk", "/old") != 30 {
		t.Error("Unexpected results: %v", store.Get("disk", "/old"))
	}

	store = NewHostStore()
	err = collectDisk("fixtures/df.darwin.txt", store)
	if err != nil {
		t.Error(err)
	}
	if store.Get("disk", "/") != 7 {
		t.Error("Unexpected results: %v", store.Get("disk", "/"))
	}

	store = NewHostStore()
	err = collectDisk("", store)
	if err != nil {
		t.Error(err)
	}
	if store.Get("disk", "/") <= 0 {
		t.Error("Expected root disk to have more than 0% usage")
	}
}
