// Safe_var.go provides thread-safe variables of any type. The value of the variable is retrieved by Get.
// The value of the variable is set with Set.
//
// Version v1.0
// Date 13 Aug 2023
//
// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock

// Import Go standard library package sync
import (
	"sync" // sync
)

// Interface holds the interface for a thread-safe variable for any type
type SafeInterface[T any] interface {
	// Get returns the value of the thread-safe variable
	Get() T
	// Set sets the value of the thread-safe variable
	Set(T)
}

// SafeVariable contains the value of the thread-safe variable and a mutex.
type SafeVariable[T any] struct {
	v  T          // value of type T
	mu sync.Mutex // Mutex
}

// Get returns the value of the thread-safe variable
func (inst *SafeVariable[T]) Get() T {
	// Lock the mutex
	inst.mu.Lock()
	// Defer unlocking the mutex
	defer inst.mu.Unlock()
	// Return the value
	return inst.v
}

// Set sets the value of the thread-safe variable to v
func (inst *SafeVariable[T]) Set(v T) {
	// Lock the mutex
	inst.mu.Lock()
	// Defer unlocking the mutex
	defer inst.mu.Unlock()
	// Set the value to v
	inst.v = v
}
