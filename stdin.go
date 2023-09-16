// Package tsmock provides an interface to test and mock Stdin based on files. It reads input
// from a file and passes it to os.Stdin. It can be configured to set the visibility
// of the input and a delay in processing each line of the input from the file. The mocked Stdin
// is executed in a go routine and can be canceled with a context.
//
// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock

// Import go standard library packages and tserr
import (
	"bufio"   // bufio
	"context" // context
	"fmt"     // fmt
	"os"      // os
	"sync"    // sync
	"time"    // time

	"github.com/thorstenrie/tserr" // tserr
)

// MockStdin contains the internal state of a mocked Stdin. It holds variables for file descriptors, a time delay, an option for visibility and an error, if any.
// It stores a context cancel function and a sync wait group. Users of the mocked Stdin are expected to use the globally exported instance tsmock.Stdin.
type MockStdin struct {
	in, r, w, o *os.File                    // input, pipe and original Stdin file descriptors
	e           SafeVariable[error]         // Error, if any
	d           SafeVariable[time.Duration] // Time delay in reading input
	v           SafeVariable[bool]          // Visibility of input
	run         SafeVariable[bool]          // True if executing, false otherwise
	set         SafeVariable[bool]          // True if pip is set, false otherwise
	cancel      context.CancelFunc          // Context cancel function
	wg          sync.WaitGroup              // Sync wait group
}

var (
	// Global mocked Stdin instance initialized to store the original os.Stdin to enable os.Stdin recovery and setting visibility of Stdin input to true.
	Stdin = newStdin()
)

// Retrieve a new mocked Stdin instance. Visibility of stdin is set to true.
func newStdin() *MockStdin {
	// Retrieve a new mocked Stdin instance and set o to the original os.Stdin
	r := &MockStdin{o: os.Stdin}
	// Set visibility of stdin to true
	r.v.Set(true)
	// Mocked stdin is not executing
	r.run.Set(false)
	// Mocked stdin is not set
	r.set.Set(false)
	// Return the new instance
	return r
}

// closePipe closes the pipe, if existing.
func (stdin *MockStdin) closePipe() {
	// Close read file descriptor, if not nil
	if stdin.r != nil {
		stdin.r.Close()
	}
	// Close write file descriptor, if not nil
	if stdin.w != nil {
		stdin.w.Close()
	}
	// Close input file descriptor, if not nil
	if stdin.in != nil {
		stdin.in.Close()
	}
	// Set the file descriptors to nil
	stdin.w, stdin.r, stdin.in = nil, nil, nil
}

// Restore restores the original os.Stdin. It cancels current execution of the mocked stdin and returns the last occurring error, if any.
func (stdin *MockStdin) Restore() error {
	// Cancel the current execution of the mocked Stdin, if execution is running
	if stdin.run.Get() {
		// Return an error if cancel function is nil
		if stdin.cancel == nil {
			return tserr.NilPtr()
		}
		// Cancel stdin execution
		stdin.cancel()
		// Set cancel function to nil
		stdin.cancel = nil
	}
	// Wait for the execution of the mocked stdin to be stopped
	stdin.wg.Wait()
	// Close existing pipe, if existing
	stdin.closePipe()
	// Restore os.Stdin to original os.Stdin
	os.Stdin = stdin.o
	// Set mocked stdin execution to false
	stdin.run.Set(false)
	// Set mocked stdin to not set
	stdin.set.Set(false)
	// Return an error, if any
	return stdin.e.Get()
}

// Delay sets a time delay d for the mocked Stdin. If d is set to a value higher than zero, each line input to the mocked Stdin will be delayed by
// d. It simulates the usual Stdin behavior to receive input with a delay from the user. It returns an error if d is lower than zero.
func (stdin *MockStdin) Delay(d time.Duration) error {
	// Return an error if d is negative
	if d < 0 {
		return tserr.Higher(&tserr.HigherArgs{Var: "d", Actual: int64(d), LowerBound: 0})
	}
	// Set time delay to d
	stdin.d.Set(d)
	// Return nil
	return nil
}

// Visibility sets the visibility of the Stdin input to v. If v is true, the simulated Stdin input is printed to Stdout, which is the usual
// behavior of a terminal. If v is false, the simulated Stdin input is not printed to Stdout, which is the usual behavior for
// a secret input of a terminal, for example a password.
func (stdin *MockStdin) Visibility(v bool) {
	// Set visibility to v
	stdin.v.Set(v)
}

// Err returns the last occurring error, if any.
func (stdin *MockStdin) Err() error {
	// Return las occurring error, if any
	return stdin.e.Get()
}

// Set sets the input of the mocked Stdin to in. If a previous mock run is still being executed, Set returns an error.
func (stdin *MockStdin) Set(in *os.File) error {
	// Return an error if in is nil
	if in == nil {
		return tserr.NilPtr()
	}
	// Return an error if mocked Stdin is executing
	if stdin.run.Get() {
		return tserr.Locked("Mocked Stdin")
	}
	// Close existing pipe, if existing
	stdin.closePipe()
	// Retrieve a new pipe
	var e error
	stdin.r, stdin.w, e = os.Pipe()
	// Return an error if retrieving a new pipe fails
	if (e != nil) || (stdin.w == nil) || (stdin.r == nil) {
		stdin.Restore()
		return tserr.NotAvailable(&tserr.NotAvailableArgs{S: "os.Pipe", Err: stdin.e.Get()})
	}
	// Set input file
	stdin.in = in
	// Set os.Stdin to pipe
	os.Stdin = stdin.r
	// Set mocked stdin to set
	stdin.set.Set(true)
	// Return nil
	return nil
}

// Run starts a new go routine to write the input from in into the mocked Stdin.
// The input can be retrieved through os.Stdin, the same as it would be user input from a terminal.
// The go routine closes and exits, when all input from in has been processed or if the context is canceled.
// To execute the delay, the Sleep function is used. If the context is canceled, the execution will stop after the Sleep function completed.
// It returns an error if the mocked Stdin is already executing.
func (stdin *MockStdin) Run(ctx context.Context) error {
	// Return an error if the mocked Stdin is already executing
	if stdin.run.Get() {
		return tserr.Locked("Mocked Stdin")
	}
	// Return an error if the mocked Stdin is not set
	if !stdin.set.Get() {
		return tserr.NotSet("Mocked Stdin")
	}
	// Add to waitgroup
	stdin.wg.Add(1)
	// Set execution to true
	stdin.run.Set(true)
	// Retrieve a child context and a cancel function
	ctx, stdin.cancel = context.WithCancel(ctx)
	// Execute mocked Stdin
	go stdin.write(ctx)
	// Return nil
	return nil
}

// write writes text from in into Stdin. It is intended to be executed in a go routine.
func (stdin *MockStdin) write(ctx context.Context) {
	// Set waitgroup to done after execution finished
	defer stdin.wg.Done()
	// Set execution to false after execution finished
	defer stdin.run.Set(false)
	// Set an error and stop execution if w is nil
	if stdin.w == nil {
		stdin.e.Set(tserr.NilPtr())
		return
	}
	// Close w after execution finished
	defer stdin.w.Close()
	// Set an error and stop execution if in is nil
	if stdin.in == nil {
		stdin.e.Set(tserr.NilPtr())
		return
	}
	// Retrieve a scanner on in
	s := bufio.NewScanner(stdin.in)
	// Set break condition to false
	br := false
	// Scan scanner on in
	for s.Scan() {
		select {
		// Set break condition to true, if context is canceled
		case <-ctx.Done():
			// Break outer loop
			br = true
		default: // Otherwise, continue
		}
		// Stop scanning if break condition is true
		if br {
			break
		}
		// Set i to retrieved text from the scanner and add a newline
		i := s.Text() + "\n"
		// Write i to Stdin
		_, err := stdin.w.WriteString(i)
		// Set an error and stop execution, if WriteString fails
		if err != nil {
			stdin.e.Set(err)
			return
		}
		// Print i if Visibility is true
		if stdin.v.Get() {
			fmt.Print(i)
		}
		// Sleep for defined delay
		time.Sleep(stdin.d.Get())
	}
}
