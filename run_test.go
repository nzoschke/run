package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRun(args, stat string) (string, string) {
	Dev = Devices{
		Stderr: &bytes.Buffer{},
		Stdout: &bytes.Buffer{},
	}

	run(args, stat)

	eb, _ := ioutil.ReadAll(Dev.Stderr)
	ob, _ := ioutil.ReadAll(Dev.Stdout)

	return string(eb), string(ob)
}

func TestEcho(t *testing.T) {
	eb, ob := testRun(`echo hello`, "")
	assert.Equal(t, `EXEC: "echo hello"
EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n\n", ob)

	eb, ob = testRun(`echo -n hello`, "")
	assert.Equal(t, `EXEC: "echo -n hello"
EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n", ob)

	eb, ob = testRun(`echo "hello"`, "")
	assert.Equal(t, `EXEC: "echo \"hello\""
EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n\n", ob)

	eb, ob = testRun(`echo 'hello'`, "")
	assert.Equal(t, `EXEC: "echo 'hello'"
EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n\n", ob)

	eb, ob = testRun(`echo 'hello\n'`, "")
	assert.Equal(t, `EXEC: "echo 'hello\\n'"
EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\\n\n\n", ob)

	eb, ob = testRun(`echo -e 'hello\n'`, "")
	assert.Equal(t, `EXEC: "echo -e 'hello\\n'"
EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n\n    \n\n", ob)

}
