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
	assert.Equal(t, "    hello\n", ob)

	eb, ob = testRun(`echo -n hello`, "")
	assert.Equal(t, `EXEC: "echo -n hello"

EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello", ob)

	eb, ob = testRun(`echo "hello"`, "")
	assert.Equal(t, `EXEC: "echo \"hello\""

EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n", ob)

	eb, ob = testRun(`echo 'hello'`, "")
	assert.Equal(t, `EXEC: "echo 'hello'"

EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n", ob)

	eb, ob = testRun(`echo 'hello\n'`, "")
	assert.Equal(t, `EXEC: "echo 'hello\\n'"

EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\\n\n", ob)

	eb, ob = testRun(`echo -e 'hello\n'`, "")
	assert.Equal(t, `EXEC: "echo -e 'hello\\n'"

EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n    \n", ob)

	eb, ob = testRun(`echo -e 'hello\nworld\n'`, "")
	assert.Equal(t, `EXEC: "echo -e 'hello\\nworld\\n'"

EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "    hello\n    world\n    \n", ob)
}

func TestTrueFalse(t *testing.T) {
	eb, ob := testRun(`true`, "Testing")
	assert.Equal(t, `STAT: Testing
EXEC: "true"

EXIT: 0
TIME: 0.0s
`, eb)
	assert.Equal(t, "", ob)

	eb, ob = testRun(`false`, "Testing")
	assert.Equal(t, `STAT: Testing
EXEC: "false"

EXIT: 1
TIME: 0.0s
FAIL: Testing
`, eb)
	assert.Equal(t, "", ob)

}
