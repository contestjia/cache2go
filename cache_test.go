package cache2go

import (
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	table := CacheTable("table")
	a := "testy test"
	table.XCache("test", 1*time.Second, a, nil)
	b, err := table.GetXCached("test")
	if err != nil || b == nil || b.Data().(string) != a {
		t.Error("Error retrieving data from cache", err)
	}
}

func TestCacheExpire(t *testing.T) {
	table := CacheTable("table")
	a := "testy test"
	table.XCache("test", 1*time.Second, a, nil)
	b, err := table.GetXCached("test")
	if err != nil || b == nil || b.Data().(string) != a {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1500 * time.Millisecond)
	b, err = table.GetXCached("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestCacheNonExpiring(t *testing.T) {
	table := CacheTable("table")
	a := "testy test"
	table.XCache("test", 0, a, nil)
	time.Sleep(500 * time.Millisecond)
	b, err := table.GetXCached("test")
	if err != nil || b == nil || b.Data().(string) != a {
		t.Error("Error retrieving data from cache", err)
	}
}

func TestCacheKeepAlive(t *testing.T) {
	table := CacheTable("table")
	a := "testy test"
	table.XCache("test", 500*time.Millisecond, a, nil)
	a = "testest test"
	table.XCache("test2", 1250*time.Millisecond, a, nil)
	b, err := table.GetXCached("test")
	if err != nil || b == nil || b.Data().(string) != "testy test" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(200 * time.Millisecond)
	b.KeepAlive()
	time.Sleep(750 * time.Millisecond)
	b, err = table.GetXCached("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
	b, err = table.GetXCached("test2")
	if err != nil || b == nil || b.Data().(string) != "testest test" {
		t.Error("Error retrieving data from cache", err)
	}
	time.Sleep(1500 * time.Millisecond)
	b, err = table.GetXCached("test2")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlush(t *testing.T) {
	table := CacheTable("table")
	a := "testy test"
	table.XCache("test", 10*time.Second, a, nil)
	time.Sleep(1000 * time.Millisecond)
	table.XFlush()
	b, err := table.GetXCached("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestFlushNoTimout(t *testing.T) {
	table := CacheTable("table")
	a := "testy test"
	table.XCache("test", 10*time.Second, a, nil)
	table.XFlush()
	b, err := table.GetXCached("test")
	if err == nil || b != nil {
		t.Error("Error expiring data")
	}
}

func TestMassive(t *testing.T) {
	table := CacheTable("table")
	val := "testy test"
	count := 100000
	for i := 0; i < count; i++ {
		key := "test_" + strconv.Itoa(i)
		table.XCache(key, 2*time.Second, val, nil)
	}
	for i := 0; i < count; i++ {
		key := "test_" + strconv.Itoa(i)
		d, err := table.GetXCached(key)
		if err != nil || d == nil || d.Data().(string) != val {
			t.Error("Error retrieving data")
		}
	}
	if table.XCacheCount() != count {
		t.Error("Data count mismatch")
	}
}
