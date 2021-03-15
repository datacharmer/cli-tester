package cli_tester

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// ExpectedFuncType is a function that can examine a result and decide whether it passes the check
type ExpectedFuncType = func(outText, errText string) error

// CliSet is the low level element for a CLI test
type CliSet struct {
	Name           string           // The name of the set
	Command        CommandSet       // The main command to execute
	ExpectedRegexp *regexp.Regexp   // Expected result as evaluated by a regular expression
	ExactExpected  string           // Exact expected result
	ExpectedFunc   ExpectedFuncType // A function capable of checking the result
	Timeout        time.Duration    // An optional timeout (for future use)
}

// CliTest is the full test to be executed
type CliTest struct {
	Executable     string      // The name of the executable to run
	Version        string      // Version of the executable
	VersionCommand string      // The command to run to get the version
	StopOnError    bool        // Halt the execution if one of the sets fails
	Sets           []*CliSet   // The tests to run
	Logger         *log.Logger // If defined, enables logging
}

// Validate will make sure that a CliTest structure contains the necessary data
func (ct *CliTest) Validate() error {
	fullExecutablePath, err := exec.LookPath(ct.Executable)
	if err != nil {
		return fmt.Errorf("[CliTest check] error searching for executable '%s': %s", ct.Executable, err)
	}
	if fullExecutablePath == "" {
		return fmt.Errorf("[CliTest check] executable '%s' not found in PATH", ct.Executable)
	}
	if ct.Version != "" {
		if ct.VersionCommand == "" {
			return fmt.Errorf("[version check] no command was provided to retrieve the version for '%s'", ct.Executable)
		}
		outText, errText, err := runCmdCtrlArgs(ct.Executable, ct.VersionCommand)
		if err != nil {
			return err
		}
		if errText != "" {
			return fmt.Errorf("[version check] expecting no stderr, got '%s'", errText)
		}
		if !strings.Contains(outText, ct.Version) {
			return fmt.Errorf("[version check] Command '%s %s' - expected version '%s' - got '%s'", ct.Executable, ct.VersionCommand, ct.Version, outText)
		}
	}
	return nil
}

// NewCliTest initializes a new CliTest
func NewCliTest(executable, version, versionCommand string, stopOnError bool, logger *log.Logger) *CliTest {
	return &CliTest{
		Executable:     executable,
		Version:        version,
		VersionCommand: versionCommand,
		StopOnError:    stopOnError,
		Sets:           nil,
		Logger:         logger,
	}
}

func (ct *CliTest) Logf(format string, v ...interface{}) {
	if ct.Logger == nil {
		return
	}
	ct.Logger.Printf(format, v...)
}

// Add adds a new set to the CLI test
func (ct *CliTest) Add(cliSet *CliSet, expectedRegexp string) error {
	if expectedRegexp != "" {
		re, err := regexp.Compile(expectedRegexp)
		if err != nil {
			return SetCliRunError(cliSet, "error compiling regex '%s': %s", expectedRegexp, err)
		}
		cliSet.ExpectedRegexp = re
	}
	if cliSet.Name == "" {
		return SetCliRunError(cliSet, "[CliTest.Add] missing 'Name' component ")
	}
	if cliSet.Command.Command == "" {
		return SetCliRunError(cliSet, "[CliTest.Add] missing 'Command' component ")
	}
	ct.Sets = append(ct.Sets, cliSet)
	ct.Logf("Added command %s", cliSet.Command)
	return nil
}

// RunSet runs a specific set
func (ct *CliTest) RunSet(cliSet *CliSet) *CliRunError {
	ct.Logf("Running command '%s': %s", cliSet.Name, cliSet.Command)
	stdOut, errOut, err := runCmdCtrlArgs(cliSet.Command.Command, cliSet.Command.Args...)
	cliSet.Command.StdOut = stdOut
	cliSet.Command.ErrOut = errOut
	if err != nil {
		ct.Logf("[RunSet] execution error %s", err)
		ct.Logf("[RunSet] stdOut <%s>", stdOut)
		ct.Logf("[RunSet] errOut <%s>", errOut)
		return SetCliRunError(cliSet, "[RunSet] execution error %s", err)
	}
	ct.Logf("[RunSet] execution concluded without errors")
	return nil
}

// Run executes all sets in the CLI test
func (ct *CliTest) Run() []*CliRunError {
	err := ct.Validate()
	if err != nil {
		return []*CliRunError{SetCliRunError(nil, "[CliTest check] %s", err)}
	}
	var errors []*CliRunError
	for i, set := range ct.Sets {
		fmt.Printf("%5d %s\n", i, set.Name)
		err := ct.RunSet(set)
		failed := false
		if err != nil {
			errors = append(errors, err)
			failed = true
		}
		if set.ExactExpected != "" && !failed {
			if set.ExactExpected != set.Command.StdOut {
				ct.Logf("[CliTest.Run] error matching exact value. Wanted '%s', got '%s'", set.ExactExpected, set.Command.StdOut)
				errors = append(errors, SetCliRunError(set, "[CliTest.Run] error matching exact value. Wanted '%s', got '%s'", set.ExactExpected, set.Command.StdOut))
				failed = true
			}
		}
		if set.ExpectedRegexp != nil && !failed {
			if !set.ExpectedRegexp.MatchString(set.Command.StdOut) {
				ct.Logf("[CliTest.Run] error matching expected '%s'", set.ExpectedRegexp.String())
				errors = append(errors, SetCliRunError(set, "[CliTest.Run] error matching expected '%s'", set.ExpectedRegexp.String()))
				failed = true
			}
		}
		if set.ExpectedFunc != nil && !failed {
			err := set.ExpectedFunc(set.Command.StdOut, set.Command.ErrOut)
			if err != nil {
				ct.Logf("[CliTest.Run] error matching expected func '%s'", err)
				errors = append(errors, SetCliRunError(set, "error matching expected func: %s", err))
				failed = true
			}
		}
		if failed && ct.StopOnError {
			return errors
		}
	}
	return errors
}
