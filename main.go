package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	shellquote "github.com/kballard/go-shellquote"
)

type PrefixWriter struct {
	w io.Writer
}

func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	b := []byte("    ")
	b = append(b, p...)
	_, err = pw.w.Write(b)
	return n, err
}

func main() {
	name := os.Args[1]
	args := os.Args[2:]

	s := fmt.Sprintf("RUN: %s '%s'\n", name, strings.Join(args, "' '"))
	os.Stderr.WriteString(s)

	sargs := []string{"-c", shellquote.Join(os.Args[1:]...)}
	cmd := exec.Command("bash", sargs...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	e, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	o, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		io.Copy(&PrefixWriter{os.Stderr}, e)
	}()

	go func() {
		io.Copy(&PrefixWriter{os.Stdout}, o)
	}()

	code := 0
	err = cmd.Wait()
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			code = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", name, args)
			code = 1
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		code = ws.ExitStatus()
	}

	s = fmt.Sprintf("CODE: %d\n", code)
	os.Stderr.WriteString(s)

	os.Exit(code)
}
