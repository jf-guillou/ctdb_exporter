package main

import (
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

	return "", fmt.Errorf("unexpected command : %v", arg)
}

func TestIsMasterNode(t *testing.T) {
	mn, err := isMasterNode(testCommandRunner)
	if err != nil {
		t.Error(err)
	}

	if !mn {
		t.Error(fmt.Errorf("expected true but got %v", mn))
	}
}

func TestScrapeStatus(t *testing.T) {
	status, err := scrapeStatus(testCommandRunner)
	if err != nil {
		t.Error(err)
	}

	if status == nil {
		t.Error("expected scrapeStatus to return something but got nothing")
	}

	expectedOutput := Status{
		"0", "0.0.0.1", 0, 0, 0, 0, 0, 0, 1, 1,
	}

	if status[0] != expectedOutput {
		t.Error("expected scraped status to be correctly parsed")
	}
}
