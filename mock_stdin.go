// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/thorstenrie/tserr"
)

type MockStdin struct {
	in, r, w, o *os.File
	e           error
	d           time.Duration
	v           bool // Visibility of input
}

var (
	Stdin = &MockStdin{o: os.Stdin, v: true}
)

func (stdin *MockStdin) Restore() error {
	if stdin.r != nil {
		stdin.e = stdin.r.Close()
	}
	if stdin.w != nil {
		stdin.e = stdin.w.Close()
	}
	if stdin.in != nil {
		stdin.e = stdin.in.Close()
	}
	stdin.w, stdin.r, stdin.in = nil, nil, nil
	os.Stdin = stdin.o
	return stdin.e
}

func (stdin *MockStdin) Delay(d time.Duration) error {
	if d < 0 {
		stdin.e = tserr.Higher(&tserr.HigherArgs{Var: "d", Actual: int64(d), LowerBound: 0})
		return stdin.e
	}
	stdin.d = d
	return nil
}

func (stdin *MockStdin) Visibility(v bool) {
	stdin.v = v
}

func (stdin *MockStdin) Err() error {
	return stdin.e
}

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
	go stdin.write()
	return nil
}

func (stdin *MockStdin) write() {
	s := bufio.NewScanner(stdin.in)
	for s.Scan() {
		i := s.Text() + "\n"
		_, err := stdin.w.WriteString(i)
		if err != nil {
			return
		}
		if stdin.v {
			fmt.Print(i)
		}
		time.Sleep(stdin.d)
	}
	stdin.w.Close()
}
