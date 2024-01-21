package pkg

import "testing"

var kv = NewStore[string, int]()

func assert[A comparable](t *testing.T, s string, comparison ...A) {
	if comparison[0] != comparison[1] {
		t.Errorf("%s Expected: %v, got %v", s, comparison[1], comparison[0])
	}
}

func TestKVsetget(t *testing.T) {
	kv.set("HasFive", 5)
	kv.set("HasOne", 1)

	val, found := kv.get("HasFive")
	assert(t, "get('HasFive')", val, 5)
	assert(t, "get('HasFive')", found, true)

	val, found = kv.get("HasOne")
	assert(t, "get('HasOne')", val, 1)
	assert(t, "get('HasOne')", found, true)

	val, found = kv.get("HasNine")
	assert(t, "get('HasNine')", val, 0)
	assert(t, "get('HasNine')", found, false)

	kv.set("HasFive", 55)
	val, found = kv.get("HasFive")
	assert(t, "get('HasFive')", val, 55)
	assert(t, "get('HasFive')", found, true)
}

func TestKVCount(t *testing.T) {
	val := kv.Count()
	assert(t, "Count()", val, 2)
}

func TestKVdelete(t *testing.T) {
	kv.delete("HasFive")
	kv.delete("HasNine")

	val := kv.Count()
	assert(t, "Count()", val, 1)
}

func TestKV(t *testing.T) {
	kv = NewStore[string, int]()
	kv.Begin()
	kv.Set("X", 200)
	kv.Set("Y", 14)

	// Test value of X
	val, found := kv.Get("X")
	assert(t, "[Session 1] Get('X')", val, 200)
	assert(t, "[Session 1] Get('X')", found, true)

	// Test value of Y
	val, found = kv.Get("Y")
	assert(t, "[Session 1] Get('Y')", val, 14)
	assert(t, "[Session 1] Get('Y')", found, true)

	// Test value of Z
	val, found = kv.Get("Z")
	assert(t, "[Session 1] Get('Z')", val, 0)
	assert(t, "[Session 1] Get('Z')", found, false)

	kv.Begin()

	// Test value of X
	val, found = kv.Get("X")
	assert(t, "[Session 2] Get('X')", val, 200)
	assert(t, "[Session 2] Get('X')", found, true)

	kv.Set("Z", 5000)

	// Test value of Z
	val, found = kv.Get("Z")
	assert(t, "[Session 2] Get('Z')", val, 5000)
	assert(t, "[Session 2] Get('Z')", found, true)

	kv.End()

	// Test value of X
	val, found = kv.Get("X")
	assert(t, "[Session 1] Get('X')", val, 200)
	assert(t, "[Session 1] Get('X')", found, true)

	// Test value of Y
	val, found = kv.Get("Y")
	assert(t, "[Session 1] Get('Y')", val, 14)
	assert(t, "[Session 1] Get('Y')", found, true)

	// Test value of Z
	val, found = kv.Get("Z")
	assert(t, "[Session 1] Get('Z')", val, 0)
	assert(t, "[Session 1] Get('Z')", found, false)

	kv.Begin()
	kv.Set("O", 1000)

	// Test value of O
	val, found = kv.Get("O")
	assert(t, "[Session 3] Get('O')", val, 1000)
	assert(t, "[Session 3] Get('O')", found, true)

	kv.Rollback()

	// Test value of O
	val, found = kv.Get("O")
	assert(t, "[Session 3] Get('O')", val, 0)
	assert(t, "[Session 3] Get('O')", found, false)

	kv.Set("Z", 5000)
	kv.Commit()

	// Test value of Z
	val, found = kv.Get("Z")
	assert(t, "[Session 1] Get('Z')", val, 5000)
	assert(t, "[Session 1] Get('Z')", found, true)
}

func TestKVComprehensive(t *testing.T) {
	kv := NewStore[string, int]()
	kv.Begin()                // Session 1 {}
	kv.Begin()                // Session 2 {}
	kv.Set("K", 5)            // Session 2 {"K": 5}
	kv.Set("O", 10)           // Session 2 {"K": 5, "O": 10}
	kv.Begin()                // Session 3 {"K": 5, "O": 10}
	kv.Set("Z", 1)            // Session 3 {"K": 5, "O": 10, "Z": 1}
	kv.Begin()                // Session 4 {"K": 5, "O": 10, "Z": 1}
	kv.Begin()                // Session 5 {"K": 5, "O": 10, "Z": 1}
	val, found := kv.Get("K") // 5
	assert(t, "[Session 5] Get('K')", found, true)
	assert(t, "[Session 5] Get('K')", val, 5)
	kv.Set("Z", 10)          // Session 5 {"K": 5, "O": 10, "Z": 10}
	kv.End()                 // Session 4 {"K": 5, "O": 10, "Z": 1}
	val, found = kv.Get("Z") // 1
	assert(t, "[Session 4] Get('Z')", found, true)
	assert(t, "[Session 4] Get('Z')", val, 1)
	kv.Set("K", 15)          // Session 4 {"K": 15, "O": 10, "Z": 1}
	kv.Commit()              // Changes comitted, Session 4 {"K": 15, "O": 10, "Z": 1}
	kv.End()                 // Session 3 {"K": 5, "O": 10, "Z": 1}
	kv.End()                 // Session 2 {"K": 5, "O": 10}
	val, found = kv.Get("K") // 5
	assert(t, "[Session 2] Get('K')", found, true)
	assert(t, "[Session 2] Get('K')", val, 5)
	kv.End()                 // Session 1 {}
	kv.End()                 // All Sessions ended
	kv.Begin()               // Session 1 {"K": 15, "O": 10, "Z": 1}
	val, found = kv.Get("K") // 15
	assert(t, "[Session 1] Get('K')", found, true)
	assert(t, "[Session 1] Get('K')", val, 15)
}
