package repositories

import (
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
)

type TTLRbTree struct {
	treeMutex  sync.Mutex
	queueMutex sync.Mutex
	cleanMutex sync.Mutex
	tree       *treemap.Map
	queue      *linkedlistqueue.Queue
}

func NewTTLRbTree() *TTLRbTree {
	return &TTLRbTree{
		tree:  treemap.NewWithStringComparator(),
		queue: linkedlistqueue.New()}
}

type queueNode struct {
	key    string
	expiry time.Time
}

// sets a key (string) and an amount of time until this key expires.
func (t *TTLRbTree) Set(key string, expiresIn time.Duration) {
	t.treeMutex.Lock()
	t.queueMutex.Lock()

	expiry := time.Now().Add(expiresIn)
	t.tree.Put(key, expiry)
	t.queue.Enqueue(queueNode{key: key, expiry: expiry})

	t.queueMutex.Unlock()
	t.treeMutex.Unlock()
}

// deletes a key (string) and returns true if the key exists and was deleted
func (t *TTLRbTree) Del(key string) bool {
	t.treeMutex.Lock()
	defer t.treeMutex.Unlock()

	currentTime := time.Now()
	fetched, found := t.tree.Get(key)
	if !found {
		return false //nothing to delete
	}

	expiry := fetched.(time.Time)
	if expiry.Before(currentTime) {
		return false //expired
	}

	t.tree.Remove(key)
	return true
}

// removes all nodes from the tree that are expired. returns number of records deleted (mostly for testing)
func (t *TTLRbTree) Clean() int {
	if !t.cleanMutex.TryLock() {
		return 0
	}
	currentTime := time.Now()
	var counter int = 0
	for {
		t.treeMutex.Lock()
		t.queueMutex.Lock()

		node, exists := t.queue.Peek()
		if !exists || node == nil || node.(queueNode).expiry.After(currentTime) {
			t.queueMutex.Unlock()
			t.treeMutex.Unlock()
			break
		}
		t.queue.Dequeue()
		t.tree.Remove(node.(queueNode).key)

		t.queueMutex.Unlock()
		t.treeMutex.Unlock()
		counter++
	}
	t.cleanMutex.Unlock()
	return counter
}
