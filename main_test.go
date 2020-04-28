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
