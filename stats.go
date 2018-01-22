package main

import (
	"reflect"
	"strings"
)

type Stats struct {
	TotalConnectionsReceived string `json:"total_connections_received"`
	TotalCommandsProcessed   string `json:"total_commands_processed"`
	InstantaneousOpsPerSec   string `json:"instantaneous_ops_per_sec"`
	TotalNetInputBytes       string `json:"total_net_input_bytes"`
	TotalNetOutputBytes      string `json:"total_net_output_bytes"`
	InstantaneousInputKbps   string `json:"instantaneous_input_kbps"`
	InstantaneousOutputKbps  string `json:"instantaneous_output_kbps"`
	RejectedConnections      string `json:"rejected_connections"`
	SyncFull                 string `json:"sync_full"`
	SyncPartialOk            string `json:"sync_partial_ok"`
	SyncPartialErr           string `json:"sync_partial_err"`
	ExpiredKeys              string `json:"expired_keys"`
	EvictedKeys              string `json:"evicted_keys"`
	KeyspaceHits             string `json:"keyspace_hits"`
	KeyspaceMisses           string `json:"keyspace_misses"`
	PubsubChannels           string `json:"pubsub_channels"`
	PubsubPatterns           string `json:"pubsub_patterns"`
	LatestForkUsec           string `json:"latest_fork_usec"`
	MigrateCachedSockets     string `json:"migrate_cached_sockets"`
	SlaveExpiresTrackedKeys  string `json:"slave_expires_tracked_keys"`
	ActiveDefragHits         string `json:"active_defrag_hits"`
	ActiveDefragMisses       string `json:"active_defrag_misses"`
	ActiveDefragKeyHits      string `json:"active_defrag_key_hits"`
	ActiveDefragKeyMisses    string `json:"active_defrag_key_misses"`
}

func setValueOnStats(stats *Stats, field string, value string) *Stats {
	v := reflect.ValueOf(stats).Elem().FieldByName(field)
	if v.IsValid() {
		v.SetString(value)
	}
	return stats
}
func (stats *Stats) Parse(data string) *Stats {
	rows := strings.Split(data, "\r\n")
	rows = rows[1:]
	//stats := &Stats{}
	for _, row := range rows {
		if row == "" {
			continue
		}
		data := strings.Split(row, ":")
		key := strings.Replace(data[0], "_", " ", -1)
		key = strings.Title(key)
		key = strings.Replace(key, " ", "", -1)
		stats = setValueOnStats(stats, key, data[1])
	}

	return stats
}

func (stats *Stats) ToMap() map[string]string {
	m := make(map[string]string)
	v := reflect.ValueOf(stats).Elem()
	typed := v.Type()
	tempintslice := []int{0}
	for i := 0; i < v.NumField(); i++ {
		tempintslice[0] = i
		name := typed.FieldByIndex(tempintslice).Name
		f := v.FieldByIndex(tempintslice)
		m[name] = f.String()
	}
	return m
}
