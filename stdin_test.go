// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock_test

// Import go standard library packages as well as tserr, tsfio and tsmock
import (
	"context" // context
	"testing" // testing
	"time"    // time

	"github.com/thorstenrie/tserr"  // tserr
	"github.com/thorstenrie/tsfio"  // tsfio
	"github.com/thorstenrie/tsmock" // tsmock
)

var (
	// Test file for stdin input
	testfile = tsfio.Filename("testdata/stdin.txt")
	// Test delay
	testdelay = 250 * time.Millisecond
)

// TestStdinV tests Stdin with the test file with visibility set to true and input delay to testdelay. It will take the contents of the
// test file as reference string and as input to stdin. If the contents received from stdin does not equal
// the contents of the test file the test fails. Also, the test fails in case of an error.
func TestStdinV(t *testing.T) {
	if e := testStdin(context.Background(), true, testdelay, t); e != nil {
		t.Error(e)
	}
}

// TestStdinV tests Stdin with the test file with visibility set to false and input delay to testdelay. It will take the contents of the
// test file as reference string and as input to stdin. If the contents received from stdin does not equal
// the contents of the test file the test fails. Also, the test fails in case of an error.
func TestStdinI(t *testing.T) {
	if e := testStdin(context.Background(), false, testdelay, t); e != nil {
		t.Error(e)
	}
}

// TestStdinTimeout tests that canceling the provided context will stop Stdin execution. The Stdin execution
// is stopped with a timeout before it is fully executed. The test fails
// if the evaulation of received Stdin and the contents of the test file equal.
func TestStdinTimeout(t *testing.T) {
	// Retrieve a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), testdelay)
	// Defer cancel function
	defer cancel()
	// The test fails if the evaluation of the received Stdin and the contects of the test file equal
	if e := testStdin(ctx, true, testdelay, t); e == nil {
		t.Error(tserr.NilFailed("testStdin"))
	}
}

// TestNegativeDelay tests if Delay returns an error in case of a negative value. The test
// fails if Delay returns nil.
func TestNegativeDelay(t *testing.T) {
	if e := tsmock.Stdin.Delay(-1); e == nil {
		t.Error(tserr.NilFailed("Delay"))
	}
}

// TestnilFile tests if Set returns an error in case of nil. The test
// fails if Set returns nil.
func TestNilFile(t *testing.T) {
	if e := tsmock.Stdin.Set(nil); e == nil {
		t.Error(tserr.NilFailed("Set"))
	}
}
