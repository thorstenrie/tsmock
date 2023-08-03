// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/thorstenrie/tserr"
	"github.com/thorstenrie/tsfio"
	"github.com/thorstenrie/tsmock"
)

var (
	testfile = tsfio.Filename("testdata/stdin.txt")
)

func TestStdinV(t *testing.T) {
	testStdin(true, t)
}

func TestStdinI(t *testing.T) {
	testStdin(false, t)
}

func testStdin(v bool, t *testing.T) {
	if t == nil {
		panic(tserr.NilPtr())
	}
	ref, err := tsfio.ReadFile(testfile)
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "ReadFile", Fn: string(testfile), Err: err}))
	}
	fs, err := tsfio.OpenFile(testfile)
	if err != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "OpenFile", Fn: string(testfile), Err: err}))
	}
	sref := string(ref)
	test := ""
	defer fs.Close()
	if e := tsmock.Stdin.Delay(time.Millisecond); e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Delay", Fn: "Stdin", Err: e}))
	}
	tsmock.Stdin.Visibility(v)
	if e := tsmock.Stdin.Set(fs); e != nil {
		t.Error(tserr.Op(&tserr.OpArgs{Op: "Set", Fn: string(testfile), Err: e}))
	}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		test += s.Text() + "\n"
	}
	if e := tsmock.Stdin.Err(); err != nil {
		t.Error(tserr.Return(&tserr.ReturnArgs{Op: "Err", Actual: fmt.Sprint(e), Want: "nil"}))
	}
	if test != sref {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: string(testfile), Want: sref, Actual: test}))
	}
	tsmock.Stdin.Restore()
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
