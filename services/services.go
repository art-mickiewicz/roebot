package services

import (
	"fmt"

	c "github.com/robfig/cron/v3"
)

var variables = make(map[string]Variable)
var services = make(map[string]func())
var cron *c.Cron

type Variable struct {
	Name        string
	Value       string
	Description string
}

func init() {
	cron = c.New()
}

func RegisterVariable(name string, description string) {
	variables[name] = Variable{Name: name, Description: description}
}

func GetVariablesValues() map[string]string {
	vals := make(map[string]string)
	for k, v := range variables {
		vals[k] = v.Value
	}
	return vals
}

func GetVariablesInfo() map[string]string {
	info := make(map[string]string)
	for k, v := range variables {
		info[k] = fmt.Sprintf("%s = %s - %s", k, v.Value, v.Description)
	}
	return info
}

func GetValue(k string) string {
	if v, ok := variables[k]; ok {
		return v.Value
	} else {
		return ""
	}
}

func SetValue(key string, val string) bool {
	if v, ok := variables[key]; ok {
		v.Value = val
		variables[key] = v
		return true
	} else {
		return false
	}
}

func RegisterService(name string, cronspec string, synchronizer func()) {
	if len(cronspec) > 0 {
		cron.AddFunc(cronspec, synchronizer)
	}
	services[name] = synchronizer
}

func SyncAll() {
	for _, srv := range services {
		srv()
	}
}

func StartCron() {
	cron.Start()
}

func StopCron() {
	cron.Stop()
}
