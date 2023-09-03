// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock_test

// Import go standard library packages as well as tserr, tsfio and tsmock
import (
	"testing" // testing
	"time"    // time

	"github.com/thorstenrie/tserr"  // tserr
	"github.com/thorstenrie/tsfio"  // tsfio
	"github.com/thorstenrie/tsmock" // tsmock
)

var (
	testfile = tsfio.Filename("testdata/stdin.txt")
)

func TestStdinV(t *testing.T) {
	testStdin(true, 250*time.Millisecond, t)
}

func TestStdinI(t *testing.T) {
	testStdin(false, 250*time.Millisecond, t)
}

func TestNegativeDelay(t *testing.T) {
	if e := tsmock.Stdin.Delay(-1); e == nil {
		t.Error(tserr.NilFailed("Delay"))
	}
}

func TestNilFile(t *testing.T) {
	if e := tsmock.Stdin.Set(nil); e == nil {
		t.Error(tserr.NilFailed("Set"))
	}
}
