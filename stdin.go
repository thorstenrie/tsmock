// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock

// Import go standard library packages and tserr
import (
	"bufio" // bufio
	"context"
	"fmt" // fmt
	"os"  // os
	"sync"
	"time" // time

	"github.com/thorstenrie/tserr" // tserr
)

// MockStdin cointains the internal state of a mocked Stdin. It holds variables for file descriptors, a time delay, an option for visibility and an error, if any.
// Users of the mocked Stdin are expected to use the globally exported instance tsmock.Stdin.
type MockStdin struct {
	in, r, w, o *os.File                    // input, pipe and original Stdin file descriptors
	e           error                       // Error, if any
	d           SafeVariable[time.Duration] // Time delay in reading input
	v           SafeVariable[bool]          // Visibility of input
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

var (
	// Global mocked Stdin instance initialized to store the original os.Stdin to enable os.Stdin recovery and setting visibility of Stdin input to true.
	Stdin = NewStdin()
)

func NewStdin() *MockStdin {
	r := &MockStdin{o: os.Stdin}
	r.v.Set(true)
	return r
}

// Restore restores the original os.Stdin. It returns the last occurring error, if any.
func (stdin *MockStdin) Restore() error {
	if stdin.cancel != nil {
		stdin.cancel()
		stdin.cancel = nil
	}
	stdin.wg.Wait()
	// Close read file descriptor, if not nil
	if stdin.r != nil {
		stdin.e = stdin.r.Close()
	}
	// Close write file descriptor, if not nil
	if stdin.w != nil {
		stdin.e = stdin.w.Close()
	}
	// Close input file descriptor, if not nil
	if stdin.in != nil {
		stdin.e = stdin.in.Close()
	}
	// Set the file descriptors to nil
	stdin.w, stdin.r, stdin.in = nil, nil, nil
	// Restore os.Stdin to original os.Stdin
	os.Stdin = stdin.o
	// Return an error, if any
	return stdin.e
}

// Delay sets a time delay d for the mocked Stdin. If d is set to a value higher than zero, each line input to the mocked Stdin will be delayed by
// d. It simulates the usual Stdin behavior to receive input with a delay from the emulated user. It returns an error if d is lower than zero.
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

// Err returns the last occurring error, if any. It is blocked until writing to the mocked Stdin is completed.
func (stdin *MockStdin) Err() error {
	// Return las occurring error, if any
	return stdin.e
}

// Set sets the input of the mocked Stdin to in. The previous mocked Stdin is closed, if any. It starts a new go routine to write the input from in into the mocked Stdin.
// The go routine closes and exits, when all input from in has been processed. The input can be
// retrieved through os.Stdin, the same as it would be user input from a terminal. It returns an error, if any.
// It is blocked until writing to the mocked Stdin is completed.
func (stdin *MockStdin) Set(in *os.File) error {
	if in == nil {
		return tserr.NilPtr()
	}
	stdin.Restore()
	stdin.r, stdin.w, stdin.e = os.Pipe()
	if (stdin.e != nil) || (stdin.w == nil) || (stdin.r == nil) {
		stdin.Restore()
		return tserr.NotAvailable(&tserr.NotAvailableArgs{S: "os.Pipe", Err: stdin.e})
	}
	stdin.in = in
	os.Stdin = stdin.r
	return nil
}

func (stdin *MockStdin) Run(ctx context.Context) {
	ctx, stdin.cancel = context.WithCancel(ctx)
	stdin.wg.Add(1)
	go stdin.write(ctx)

}

func (stdin *MockStdin) write(ctx context.Context) {
	defer stdin.w.Close() // Todo: Add error handling
	defer stdin.wg.Done()
	s := bufio.NewScanner(stdin.in)
	br := false
	for s.Scan() {
		select {
		case <-ctx.Done():
			// Break outer loop
			br = true
		default:
		}
		if br {
			break
		}
		i := s.Text() + "\n"
		_, err := stdin.w.WriteString(i)
		if err != nil {
			stdin.e = err
			return
		}
		if stdin.v.Get() {
			fmt.Print(i)
		}
		time.Sleep(stdin.d.Get())
	}
}
