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
	testfile = tsfio.Filename("stdin.txt")
	// Test file contents
	contents = "Aragorn\nBoromir\nGandalf\nGimli\nLegolas\n"
	// Test delay
	testdelay = 250 * time.Millisecond
)

// TODO: Set again, Run without Set

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
// if the evaluation of received Stdin and the contents of the test file equal.
func TestStdinTimeout(t *testing.T) {
	// Retrieve a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), testdelay)
	// Defer cancel function
	defer cancel()
	// The test fails if the evaluation of the received Stdin and the contents of the test file equal
	if e := testStdin(ctx, true, testdelay, t); e == nil {
		t.Error(tserr.NilFailed("testStdin"))
	}
}

// TestStdinRestore tests Restore to cancel a running execution of a mocked Stdin. The test fails
// if Restore returns an error of Stdin has an error in Err.
func TestStdinRestore(t *testing.T) {
	// Read reference data and open a stdin test input file for testing.
	// Visibility of the input is set to v. The input delay is set to d.
	// The input file of stdin is set to the stdin test input file
	_, fs := testStdinSetup(true, testdelay, t)
	// Defer closing the retrieved file.
	defer fs.Close()
	// Mock Stdin
	if e := tsmock.Stdin.Run(context.Background()); e != nil {
		// The test fails if Run returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Run", Fn: "Stdin", Err: e}))
	}
	if e := tsmock.Stdin.Restore(); e != nil {
		// The test fails if Restore returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Restore", Fn: "Stdin", Err: e}))
	}
	// Restore Stdin. The test fails if Stdin has an error in Err.
	testStdinClose(t)
}

// TestStdinSet tests Set to return an error, when used while a mocked Stdin is executing.
// The test fails if Set returns nil or Stdin has an error in Err.
func TestStdinSet(t *testing.T) {
	// Read reference data and open a stdin test input file for testing.
	// Visibility of the input is set to v. The input delay is set to d.
	// The input file of stdin is set to the stdin test input file
	_, fs := testStdinSetup(true, testdelay, t)
	// Defer closing the retrieved file.
	defer fs.Close()
	// Mock Stdin
	if e := tsmock.Stdin.Run(context.Background()); e != nil {
		// The test fails if Run returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Run", Fn: "Stdin", Err: e}))
	}
	if e := tsmock.Stdin.Set(fs); e == nil {
		// The test fails if Set returns nil
		t.Error(tserr.NilFailed("Set"))
	}
	// Restore Stdin. The test fails if Stdin has an error in Err.
	testStdinClose(t)
}

// TestStdinSetAgain tests setting mocked stdin again. The test fails if Set returns an error or if any other
// error occurs during the execution.
func TestStdinSetAgain(t *testing.T) {
	// Read reference data and open a stdin test input file for testing.
	// Visibility of the input is set to v. The input delay is set to d.
	// The input file of stdin is set to the stdin test input file
	_, fs := testStdinSetup(true, testdelay, t)
	// Defer closing the retrieved file.
	defer fs.Close()
	// Set again
	if e := tsmock.Stdin.Set(fs); e != nil {
		// The test fails if Set returns nil
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set", Fn: "Mocked Stdin", Err: e}))
	}
	// Mock Stdin
	if e := tsmock.Stdin.Run(context.Background()); e != nil {
		// The test fails if Run returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Run", Fn: "Stdin", Err: e}))
	}
	// Restore Stdin. The test fails if Stdin has an error in Err.
	testStdinClose(t)
}

// TestStdinRunWithoutSet tests Run to return an error if Run is called again while the previous mocked Stdin execution has not finished yet. The test
// fails if Run returns nil. The test fails if any other error occurs.
func TestStdinRunAgain(t *testing.T) {
	// Read reference data and open a stdin test input file for testing.
	// Visibility of the input is set to v. The input delay is set to d.
	// The input file of stdin is set to the stdin test input file
	_, fs := testStdinSetup(true, testdelay, t)
	// Defer closing the retrieved file.
	defer fs.Close()
	// Set again
	if e := tsmock.Stdin.Set(fs); e != nil {
		// The test fails if Set returns nil
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set", Fn: "Mocked Stdin", Err: e}))
	}
	// Mock Stdin
	if e := tsmock.Stdin.Run(context.Background()); e != nil {
		// The test fails if Run returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Run", Fn: "Stdin", Err: e}))
	}
	// Mock Stdin again and expect an error
	if e := tsmock.Stdin.Run(context.Background()); e == nil {
		// The test fails if Run returns an error
		t.Error(tserr.NilFailed("Run"))
	}
	// Restore Stdin. The test fails if Stdin has an error in Err.
	testStdinClose(t)
}

// TestStdinRunWithoutSet tests Run to return an error if the mocked Stdin was not set before. The test
// fails if Run returns nil.
func TestStdinRunWithoutSet(t *testing.T) {
	// Mock Stdin
	if e := tsmock.Stdin.Run(context.Background()); e == nil {
		// The test fails if Run returns an error
		t.Error(tserr.NilFailed("Run"))
	}
}

// TestNegativeDelay tests if Delay returns an error in case of a negative value. The test
// fails if Delay returns nil.
func TestNegativeDelay(t *testing.T) {
	if e := tsmock.Stdin.Delay(-1); e == nil {
		t.Error(tserr.NilFailed("Delay"))
	}
}

// TestNilFile tests if Set returns an error in case of nil. The test
// fails if Set returns nil.
func TestNilFile(t *testing.T) {
	if e := tsmock.Stdin.Set(nil); e == nil {
		t.Error(tserr.NilFailed("Set"))
	}
}
