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

	return n, nil
}

func write(w io.Writer, line []byte) error {
	if line == nil {
		return nil
	}

	pre := []byte("    ")
	b := append(pre, line...)
	b = append(b, []byte("\n")...)
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

	if *s != "" {
		fmt.Fprintf(os.Stderr, "STAT: %s\n", *s)
	}
	fmt.Fprintf(os.Stderr, "EXEC: %q\n", args)

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
			fmt.Fprintf(os.Stderr, "Could not get exit code for failed program: %s", args)
			code = 1
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		code = ws.ExitStatus()
	}

	dur := time.Since(now)

	fmt.Fprintf(os.Stderr, "EXIT: %d\n", code)
	fmt.Fprintf(os.Stderr, "TIME: %0.1fs\n", dur.Seconds())

	if code != 0 && *s != "" {
		fmt.Fprintf(os.Stderr, "STAT: %s failed\n", *s)
	}

	os.Exit(code)
}
