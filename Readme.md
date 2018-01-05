# run

`run` is a utility for running commands in test environments. It adds human and machine friendly execution info to stderr, while preserving stdin, stdout and exit codes for the command.

The output is compatible with the [BIOS GitHub App](https://www.mixable.net/docs/bios/).

```console
$ go install github.com/nzoschke/run

$ run -h
usage: run [options] cmd [args]
  -help
    	show usage
  -s string
    	Add status messages to stderr

$ run true
EXEC: "true"
EXIT: 0
TIME: 0.0s

$ run false
EXEC: "false"
EXIT: 1
TIME: 0.0s

$ run -s Cloning git clone https://github.com/nzoschke/run
STAT: Cloning
EXEC: "git clone https://github.com/nzoschke/run"
    Cloning into 'run'...
EXIT: 0
TIME: 0.4s

$ run -s Cloning git clone https://github.com/nzoschke/run
STAT: Cloning
EXEC: "git clone https://github.com/nzoschke/run"
    fatal: destination path 'run' already exists and is not an empty directory.
EXIT: 128
TIME: 0.0s
STAT: Cloning failed
```
