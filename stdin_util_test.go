// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock_test

// Import go standard library packages as well as tserr, tsfio and tsmock
import (
	"bufio"   // bufio
	"context" // context

	// fmt
	"os"      // os
	"testing" // testing
	"time"    // time

	"github.com/thorstenrie/tserr"  // tserr
	"github.com/thorstenrie/tsfio"  // tsfio
	"github.com/thorstenrie/tsmock" // tsmock
)

// testStdin tests Stdin with the test file with visibility set to v and input delay to d. It will take the contents of the
// test file as reference string and as input to stdin. If the contents received from stdin does not equal
// the contents of the test file the test returns an error. Also, the test fails in case of an error.
func testStdin(ctx context.Context, v bool, d time.Duration, t *testing.T) error {
	// Panic if t is nil
	if t == nil {
		panic(tserr.NilPtr())
	}
	// Read reference data and open a stdin test input file for testing.
	// Visibility of the input is set to v. The input delay is set to d.
	// The input file of stdin is set to the stdin test input file
	ref, fs := testStdinSetup(v, d, t)
	// Defer closing the retrieved file.
	defer fs.Close()
	// Defer restoring Stdin. The test fails if Stdin has an error in Err.
	defer testStdinClose(t)
	// Mock Stdin
	if e := tsmock.Stdin.Run(ctx); e != nil {
		// The test returns an error if Run fails
		return e
	}
	// Scan stdin and compare the retrieved text with the reference string ref.
	if e := testStdinEval(ref, t); e != nil {
		// The test returns an error if the retrieved text does not equal the reference string ref.
		return e
	}
	return nil
}

// testStdinSetup reads reference data and opens a stdin test input file for testing. Visibility of
// the input is set with v and the input delay with d as well as the input file is set to the stdin test input file.
// It returns the reference data as string and
// the stdin test input file as *os.File. The test fails in case of an error.
func testStdinSetup(v bool, d time.Duration, t *testing.T) (string, *os.File) {
	// Panic if t is nil
	if t == nil {
		panic(tserr.NilPtr())
	}
	// Write the contents of the testfile to the testfile
	tsfio.WriteSingleStr(testfile, contents)
	// Retrieve reference data from test file contents
	ref := contents
	// Open the testfile
	fs, err := tsfio.OpenFile(testfile)
	// The test fails if OpenFile returns an error
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "OpenFile", Fn: string(testfile), Err: err}))
	}
	// Set stdin to fs
	if e := tsmock.Stdin.Set(fs); e != nil {
		// The test fails if Set returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set", Fn: string(testfile), Err: e}))
	}
	// Set visibility of stdin to v
	tsmock.Stdin.Visibility(v)
	// Set input delay to d
	if e := tsmock.Stdin.Delay(d); e != nil {
		// The test fails if Delay returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Delay", Fn: "Stdin", Err: e}))
	}
	// Return the reference string and fs
	return string(ref), fs
}

// testStdinEval scans stdin and compares retrieved text with the reference string ref.
// The test fails if the retrieved text does not equal the reference string ref.
func testStdinEval(ref string, t *testing.T) error {
	// Panic if t is nil
	if t == nil {
		panic(tserr.NilPtr())
	}
	// Initialize retrieved text with an empty string
	test := ""
	// Retrieve a new scanner on Stdin
	s := bufio.NewScanner(os.Stdin)
	// Scan stdin and add retrieved text in new lines
	for s.Scan() {
		test += s.Text() + "\n"
	}
	// The test fails if the retrieved text does not equal to the reference string ref
	if tsfio.NormNewlinesStr(test) != tsfio.NormNewlinesStr(ref) {
		return tserr.EqualStr(&tserr.EqualStrArgs{Var: string(testfile), Want: ref, Actual: test})
	}
	return nil
}

// testStdinClose restores the stdin. The test fails if Stdin has an error in Err.
func testStdinClose(t *testing.T) {
	// Panic if t is nil
	if t == nil {
		panic(tserr.NilPtr())
	}
	// The test fails if Stdin has an error in Err
	if e := tsmock.Stdin.Err(); e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Err", Fn: "Mocked Stdin", Err: e}))
	}
	// Restore Stdin
	if e := tsmock.Stdin.Restore(); e != nil {
		// The test fails if Restore returns an error
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Restore", Fn: "tsmock.Stdin", Err: e}))
	}
	// Remove testfile
	if e := tsfio.RemoveFile(testfile); e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Remove", Fn: string(testfile), Err: e}))
	}
}
