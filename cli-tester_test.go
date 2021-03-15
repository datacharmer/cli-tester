package cli_tester

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func C(args ...string) CommandSet {
	return NewCommandSet(args...)
}
func CS(args string) CommandSet {
	return NewStringCommandSet(args)
}

// TODO: add run with an expected failure
func TestCliTest_Run(t *testing.T) {
	var logger *log.Logger
	if os.Getenv("LOG_ON_SCREEN") != "" {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	} else {
		var logfileName = "test.log"
		logFile, err := os.OpenFile(logfileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			t.Errorf("error opening log file %s: %s", logfileName, err)
			return
		}
		defer logFile.Close()
		logger = log.New(logFile, "", log.Ldate|log.Ltime)
		defer os.Remove(logfileName)
	}
	ct := NewCliTest("dbdeployer", "1.59.0", "--version", true, logger)

	ct.Add(&CliSet{
		Name: "deploy single",
		//Command: C("dbdeployer", "deploy", "single", "8.0.23"),
		Command: CS("dbdeployer deploy single 8.0.23"),
	}, "started")
	ct.Add(&CliSet{
		Name: "delete single",
		//Command: C("dbdeployer", "delete", "msb_8_0_23"),
		Command: CS("dbdeployer delete msb_8_0_23"),
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
		if err.RunSet != nil {
			fmt.Printf("StdOut %#v\n", err.RunSet.Command.StdOut)
			fmt.Printf("ErrOut %#v\n", err.RunSet.Command.ErrOut)
		}
	}
}
