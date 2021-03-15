package cli_tester

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// CommandSet contains the elements needed to run a command
type CommandSet struct {
	Command string   // The name of the executable (or shell command)
	Args    []string // arguments to the command
	StdOut  string   // regular output returned from the command execution
	ErrOut  string   // output to stdErr returned from the command execution
}

// NewStringCommandSet creates a CommandSet from a string
func NewStringCommandSet(cmd string) CommandSet {
	args := strings.Split(cmd, " ")
	return NewCommandSet(args...)
}

// NewCommandSet creates a CommandSet from an array of strings
func NewCommandSet(args ...string) CommandSet {
	if len(args) == 0 {
		return CommandSet{}
	}
	return CommandSet{
		Command: args[0],
		Args:    args[1:],
	}
}

func (cs CommandSet) String() string {
	if len(cs.Args) > 0 {
		return fmt.Sprintf("%s %s", cs.Command, strings.Join(cs.Args, " "))
	}
	return cs.Command
}

// runCmdCtrlArgs executes a command, returning stdOut, stdErr, and error
func runCmdCtrlArgs(c string, args ...string) (string, string, error) {
	cmd := exec.Command(c, args...) // #nosec G204
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", "", err
	}

	slurpErr, _ := io.ReadAll(stderr)
	slurpOut, _ := io.ReadAll(stdout)
	err = cmd.Wait()

	return string(slurpOut), string(slurpErr), err
}
