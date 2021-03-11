package cli_tester

import (
	"fmt"
	"testing"
)

func C(args ...string) CommandSet {
	return NewCommandSet(args...)
}

func TestCliTest_Run(t *testing.T) {
	ct := NewCliTest("dbdeployer", "1.59.0", "--version", true)

	ct.Add(&CliSet{
		Name:    "deploy single",
		Command: C("dbdeployer", "deploy", "single", "8.0.23"),
	}, "started")
	ct.Add(&CliSet{
		Name:    "delete single",
		Command: C("dbdeployer", "delete", "msb_8_0_23"),
	}, "deleted")
	ct.Add(&CliSet{
		Name:    "deploy replication",
		Command: C("dbdeployer", "deploy", "replication", "8.0.23", "--concurrent"),
	}, "installed")
	ct.Add(&CliSet{
		Name:    "delete replication",
		Command: C("dbdeployer", "delete", "rsandbox_8_0_23"),
		//Command: C("dbdeployer", "delete", "xxxx"),
	}, "deleted")
	errors := ct.Run()
	if IsFailed(errors) {
		t.Errorf("error found during execution: %s", ErrorMessages(errors))
	}
	for i, err := range errors {
		fmt.Printf("ErrorMessage %d %#v\n", i, err.ErrorMessage)
		fmt.Printf("StdOut %#v\n", err.RunSet.Command.StdOut)
		fmt.Printf("ErrOut %#v\n", err.RunSet.Command.ErrOut)
	}
}
