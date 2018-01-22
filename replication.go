package main

import (
	"reflect"
	"strings"
)

type Replication struct {
	Role                       string `json:"role"`
	ConnectedSlaves            string `json:"connected_slaves"`
	MasterReplid               string `json:"master_replid"`
	MasterReplid2              string `json:"master_replid2"`
	MasterReplOffset           string `json:"master_repl_offset"`
	SecondReplOffset           string `json:"second_repl_offset"`
	ReplBacklogActive          string `json:"repl_backlog_active"`
	ReplBacklogSize            string `json:"repl_backlog_size"`
	ReplBacklogFirstByteOffset string `json:"repl_backlog_first_byte_offset"`
	ReplBacklogHistlen         string `json:"repl_backlog_histlen"`
}

func setValueOnReplication(stats *Replication, field string, value string) *Replication {
	v := reflect.ValueOf(stats).Elem().FieldByName(field)
	if v.IsValid() {
		v.SetString(value)
	}
	return stats
}
func (replication *Replication) Parse(data string) *Replication {
	rows := strings.Split(data, "\r\n")
	rows = rows[1:]

	for _, row := range rows {
		if row == "" {
			continue
		}
		data := strings.Split(row, ":")
		key := strings.Replace(data[0], "_", " ", -1)
		key = strings.Title(key)
		key = strings.Replace(key, " ", "", -1)
		replication = setValueOnReplication(replication, key, data[1])
	}

	return replication
}

func (replication *Replication) ToMap() map[string]string {
	m := make(map[string]string)
	v := reflect.ValueOf(replication).Elem()
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
