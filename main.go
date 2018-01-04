package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	shellquote "github.com/kballard/go-shellquote"
)

// Devices holds writers for stderr and stdout
type Devices struct {
	Stderr io.ReadWriter
	Stdout io.ReadWriter
}

// Dev is devices for normal execution
var Dev = Devices{
	Stderr: os.Stderr,
	Stdout: os.Stdout,
}

// PrefixWriter writes lines indented by 4 spaces
type PrefixWriter struct {
	w io.Writer
}

func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	buf := bytes.NewBuffer(p)
	for {
		l, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				if err := write(pw.w, l); err != nil {
					return n, err
				}
				break
			}
			return n, err
		}

		if err := write(pw.w, l); err != nil {
			return n, err
		}
	}

	return n, err
}

func write(w io.Writer, line []byte) error {
	if line == nil {
		return nil
	}

	pre := []byte("    ")
	b := append(pre, line...)
	_, err := w.Write(b)
	return err
}

func main() {
	flag.Bool("help", false, "show usage")
	s := flag.String("s", "", "Add status messages to stderr")

	flag.Usage = func() {
		fmt.Printf("usage: %s [options] cmd [args]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	args := shellquote.Join(flag.Args()...)

	if args == "" {
		flag.Usage()
		os.Exit(1)
	}

	code := run(args, *s)
	os.Exit(code)
}

func run(args, stat string) int {
	if stat != "" {
		fmt.Fprintf(Dev.Stderr, "STAT: %s\n", stat)
	}
	fmt.Fprintf(Dev.Stderr, "EXEC: %q\n", args)

	now := time.Now()
	sargs := []string{"-c", args}
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

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		io.Copy(&PrefixWriter{Dev.Stderr}, e)
	}()

	go func() {
		io.Copy(&PrefixWriter{Dev.Stdout}, o)
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
			fmt.Fprintf(Dev.Stderr, "Could not get exit code for failed program: %s", args)
			code = 1
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		code = ws.ExitStatus()
	}

	dur := time.Since(now)

	fmt.Fprintf(Dev.Stderr, "\nEXIT: %d\n", code)
	fmt.Fprintf(Dev.Stderr, "TIME: %0.1fs\n", dur.Seconds())

	if code != 0 && stat != "" {
		fmt.Fprintf(Dev.Stderr, "STAT: %s failed\n", stat)
	}

	return code
}
