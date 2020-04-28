package main

import (
	"fmt"
	"testing"
)

func testStatusCommandRunner(arg ...string) (string, error) {
	switch arg[0] {
	case "status -Y":
		return "|Node|IP|Disconnected|Banned|Disabled|Unhealthy|Stopped|Inactive|PartiallyOnline|ThisNode|\n|0|0.0.0.1|0|0|0|0|0|0|1|Y|", nil
	}

	return "", fmt.Errorf("unexpected command : %v", arg)
}

func TestScrapeStatus(t *testing.T) {
	status, err := scrapeStatus(testStatusCommandRunner)
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
