// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tsmock_test

import (
	"bufio"
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

func TestStdin(t *testing.T) {
	fs, _ := tsfio.OpenFile(testfile)
	ref, _ := tsfio.ReadFile(testfile)
	sref := string(ref)
	test := ""
	defer fs.Close()
	tsmock.Stdin.Delay(time.Second)
	tsmock.Stdin.Visibility(false)
	tsmock.Stdin.Set(fs)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		test += s.Text() + "\n"
	}
	if test != sref {
		t.Error(tserr.EqualStr(&tserr.EqualStrArgs{Var: string(testfile), Want: sref, Actual: test}))
	}
	tsmock.Stdin.Restore()
}
