package shell

import (
	"io"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Cmd struct {
	exec.Cmd

	DryRun bool
}

func (o *Cmd) Run() error {
	// Log what and where we're running
	log.WithField("dir", o.Cmd.Dir).Debugf("Running: %s", strings.Join(o.Cmd.Args, " "))

	// If dry run, don't run the command
	if o.DryRun {
		return nil
	}

	// Build writer slice
	writers := []io.Writer{
		log.StandardLogger().WriterLevel(log.TraceLevel),
	}

	// Add log writer to stdout
	writersStdout := writers
	if o.Cmd.Stdout != nil {
		writersStdout = append(writers, o.Cmd.Stdout)
	}
	o.Cmd.Stdout = io.MultiWriter(writersStdout...)

	// Add log writer to stderr
	writersStderr := writers
	if o.Cmd.Stderr != nil {
		writersStderr = append(writers, o.Cmd.Stderr)
	}
	o.Cmd.Stderr = io.MultiWriter(writersStderr...)

	return o.Cmd.Run()
}
