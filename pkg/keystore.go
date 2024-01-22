package pkg

import (
	"container/list"
)

/* Questions/Notes:
- Each session can have child sessions
- Commiting changes in the child session
  should persist those changes to the main map
	and the parent session
*/

// Type assertions are used because the list.List.Value
// stores pointers to any type, type casting tells it
// what is the stored type (or the type pointed to by the
// pointer)
// https://stackoverflow.com/a/28184823/768020

// KV is a map type
type KV[K comparable, V any] map[K]V

// Store is the Key Value store that offers
// transactional capabilities
type Store[K comparable, V any] struct {
	Map      KV[K, V]
	Sessions list.List
	Session  KV[K, V]
}

// NewStore returns a new Store
func NewStore[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		Map:      KV[K, V]{},
		Sessions: *list.New(),
	}
}

// The following non-exported methods are not indended
// for use by clients, they directly interact with the
// main map

// set is a non-exported method that directly
// sets the key values on the main map
func (Store *Store[K, V]) set(k K, v V) {
	Store.Map[k] = v
}

// get is a non-exported method that directly
// gets the values from the main map
func (Store *Store[K, V]) get(k K) (v V, found bool) {
	v, found = Store.Map[k]
	return
}

// delete is a non-exported method that directly
// deletes a key from the main map
func (Store *Store[K, V]) delete(k K) {
	_, found := Store.Map[k]
	if found {
		delete(Store.Map, k)
	}
}

// Count returns the count of keys in the main
// map
func (Store *Store[K, V]) Count() int {
	return len(Store.Map)
}

// forkSession creates a new session on a
// Begin command; If this is the first begin
// the session is initialized with the main map
// otherwise from the current ongoing session
func (Store *Store[K, V]) forkSession() {
	s := KV[K, V]{}
	if Store.Session == nil {
		for key, val := range Store.Map {
			s[key] = val
		}
	} else {
		for key, val := range Store.Session {
			s[key] = val
		}
		Store.Sessions.PushFront(Store.Session)
	}
	Store.Session = s
}

// Begin starts a new session
// A new Store does not have any session
// on-going
func (Store *Store[K, V]) Begin() {
	Store.forkSession()
}

// End ends the current session, all un-commited
// changes are lost
// After all sessions end, i.e. equal number
// of 'Ends' have been called as 'Begin', the
// session becomes nil
func (Store *Store[K, V]) End() {
	if Store.Sessions.Len() == 0 {
		Store.Session = nil
	} else {
		s := Store.Sessions.Front()
		Store.Sessions.Remove(s)
		Store.Session = s.Value.(KV[K, V])
	}
}

// Rollback discards the changes made in the
// current session, but keeps the current session
// active - which is the difference with End
func (Store *Store[K, V]) Rollback() {
	if Store.Session == nil {
		return
	}
	s := Store.Sessions.Front()
	Store.Sessions.Remove(s)
	Store.Session = KV[K, V]{}
	for key, val := range s.Value.(KV[K, V]) {
		Store.Session[key] = val
	}
	Store.Sessions.PushFront(s)
}

// Commit persists the changes in the current
// session to the main map KV, without ending the
// current session
// To-Do:
func (Store *Store[K, V]) Commit() {
	if Store.Session == nil {
		return
	}
	Store.Map = KV[K, V]{}
	for key, val := range Store.Session {
		Store.Map[key] = val
	}

	// The changes in the current session should also persist
	// to parent session
	if parent := Store.Sessions.Front(); parent != nil {
		s := KV[K, V]{}
		for key, val := range Store.Session {
			s[key] = val
		}
		Store.Sessions.Remove(parent)
		Store.Sessions.PushFront(s)
	}
}

// Set sets a particular key to a value
func (Store *Store[K, V]) Set(k K, v V) {
	if Store.Session == nil {
		return
	}
	Store.Session[k] = v
}

func (Store *Store[K, V]) Get(k K) (v V, found bool) {
	if Store.Session == nil {
		return
	}
	v, found = Store.Session[k]
	return
}

func (Store *Store[K, V]) Delete(k K) {
	if Store.Session == nil {
		return
	}
	delete(Store.Session, k)
}
