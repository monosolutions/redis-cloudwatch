package main

import (
	"log"
	"reflect"
	"strings"
)

type Cpu struct {
	UsedCpuSys          string
	UsedCpuUser         string
	UsedCpuSysChildren  string
	UsedCpuUserChildren string
}

func setValueOnCpu(cpu *Cpu, field string, value string) *Cpu {
	v := reflect.ValueOf(cpu).Elem().FieldByName(field)
	if v.IsValid() {
		v.SetString(value)
	}
	return cpu
}
func (cpu *Cpu) Parse(data string) *Cpu {
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
		log.Printf("Key %v : Value: %v", key, data[1])
		cpu = setValueOnCpu(cpu, key, data[1])
	}

	return cpu
}

func (cpu *Cpu) ToMap() map[string]string {
	m := make(map[string]string)
	v := reflect.ValueOf(cpu).Elem()
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
