package main

import (
	"fmt"
	"testing"
)

func testStatisticsCommandRunner(arg ...string) (string, error) {
	switch arg[0] {
	case "statistics -Y":
		return "CTDB version|Current time of statistics|Statistics collected since|num_clients|frozen|recovering|num_recoveries|client_packets_sent|client_packets_recv|node_packets_sent|node_packets_recv|keepalive_packets_sent|keepalive_packets_recv|node.req_call|node.reply_call|node.req_dmaster|node.reply_dmaster|node.reply_error|node.req_message|node.req_control|node.reply_control|node.req_tunnel|client.req_call|client.req_message|client.req_control|client.req_tunnel|timeouts.call|timeouts.control|timeouts.traverse|locks.num_calls|locks.num_current|locks.num_pending|locks.num_failed|total_calls|pending_calls|childwrite_calls|pending_childwrite_calls|memory_used|max_hop_count|total_ro_delegations|total_ro_revokes|num_reclock_ctdbd_latency|min_reclock_ctdbd_latency|avg_reclock_ctdbd_latency|max_reclock_ctdbd_latency|num_reclock_recd_latency|min_reclock_recd_latency|avg_reclock_recd_latency|max_reclock_recd_latency|num_call_latency|min_call_latency|avg_call_latency|max_call_latency|num_lockwait_latency|min_lockwait_latency|avg_lockwait_latency|max_lockwait_latency|num_childwrite_latency|min_childwrite_latency|avg_childwrite_latency|max_childwrite_latency|\n1|1588091528|1588085478|46|0|0|5|400051|459095|857734|353214|3620|3620|125059|0|43553|82128|0|26998|331404|217203|0|155646|9847|294440|0|0|0|1|905|0|0|0|155646|0|0|0|212806|2|0|0|5|0.004952|0.011306|0.030810|1|0.004096|0.004096|0.004096|154808|0.000004|0.001244|0.438651|905|0.002878|0.004589|0.036985|0|0.000000|0.000000|0.000000|", nil
	}

	return "", fmt.Errorf("unexpected command : %v", arg)
}

func TestScrapeStatistics(t *testing.T) {
	statistics, err := scrapeStatistics(testStatisticsCommandRunner)
	if err != nil {
		t.Error(err)
	}

	if statistics == nil {
		t.Error("expected scrapeStatistics to return something but got nothing")
	}

	expectedOutput := Statistics{
		numClients:         46,
		numRecoveries:      5,
		clientPacketsSent:  400051,
		clientPacketsRecv:  459095,
		maxHopCount:        2,
		numCallLatency:     154808,
		minCallLatency:     0.000004,
		avgCallLatency:     0.001244,
		maxCallLatency:     0.438651,
		numLockwaitLatency: 905,
		minLockwaitLatency: 0.002878,
		avgLockwaitLatency: 0.004589,
		maxLockwaitLatency: 0.036985,
	}

	if *statistics != expectedOutput {
		t.Error("expected scraped statistics to be correctly parsed")
	}
}
