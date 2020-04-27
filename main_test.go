package main

import (
	"errors"
	"fmt"
	"testing"
)

func testCommandRunner(arg ...string) (string, error) {
	switch arg[0] {
	case "pnn":
		return "0", nil
	case "recmaster":
		return "0", nil
	case "status -Y":
		return "|Node|IP|Disconnected|Banned|Disabled|Unhealthy|Stopped|Inactive|PartiallyOnline|ThisNode|\n|0|0.0.0.1|0|0|0|0|0|0|1|Y|", nil
	}

	return "", errors.New(fmt.Sprintf("unexpected command : %s", arg[0]))
}

func TestIsMasterNode(t *testing.T) {
	mn, err := isMasterNode(testCommandRunner)
	if err != nil {
		t.Error(err)
	}

	if !mn {
		t.Error(fmt.Sprintf("expected true but got %v", mn))
	}
}

func TestScrapeStatus(t *testing.T) {
	status, err := scrapeStatus(testCommandRunner)
	if err != nil {
		t.Error(err)
	}

	if status == nil {
		t.Error("expected something but got nothing")
	}

	expectedOutput := Status{
		"0", "0.0.0.1", 0, 0, 0, 0, 0, 0, 1, 1,
	}

	if status[0] != expectedOutput {
		t.Error("expected status to be correctly parsed")
	}
}
