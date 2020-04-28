package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "ctdb"
)

var (
	addr     = flag.String("web.listen-address", ":9727", "The address to listen on for HTTP requests.")
	endpoint = flag.String("web.endpoint", "/metrics", "Path under which to expose metrics.")
	ctdbBin  = flag.String("ctdb.bin-path", "/usr/bin/ctdb", "Full path to CTDB binary.")
	ctdbSudo = flag.Bool("ctdb.sudo", true, "Prefix ctdb commands with sudo.")
	pnn = ""
	recmaster = ""
)

type runner func(...string) (string, error)

func runCmd(arg ...string) (string, error) {
	cmd := exec.Command(*ctdbBin, arg...)
	if *ctdbSudo {
		// This monstrosity of a command tries to run /bin/sh -c /usr/bin/sudo /usr/bin/ctdb with proper escaping
		cmd = exec.Command("/bin/sh", append([]string{"-c"},
			strings.Join(append([]string{"/usr/bin/sudo", *ctdbBin}, arg...), " "))...)
	}
	result, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command '%v' failed with err : %v (%v)", cmd.String(), err, strings.TrimSpace(string(result)))
	}

	return strings.TrimSpace(string(result)), nil
}

func isMasterNode(run runner) (bool, error) {
	var err error
	pnn, err = run("pnn")
	if err != nil {
		return false, err
	}

	recmaster, err = run("recmaster")
	if err != nil {
		return false, err
	}

	return pnn == recmaster, nil
}

func main() {
	flag.Parse()

	prometheus.MustRegister(NewStatusCollector())
	prometheus.MustRegister(NewStatisticsCollector())

	http.Handle(*endpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
