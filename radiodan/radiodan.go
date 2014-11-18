package radiodan

import (
	"fmt"
	"time"
)

type DebugInfo struct {
	Timestamp          time.Time
	Hostname           string
	Addresses          []string
	InternetConnection bool
	Applications       []RadiodanApplication
}

type RadiodanApplication struct {
	Name      string
	LogTail   string
	Pid       int64
	IsRunning bool
	Deploy    Deploy
}

type Deploy struct {
	Name   string
	Ref    string
	Commit string
}

func (r RadiodanApplication) DeployFile() (path string) {
	path = fmt.Sprintf("/opt/radiodan/apps/%s/current/.deploy", r.Name)
	return
}

func (r RadiodanApplication) LogFile() (path string) {
	path = fmt.Sprintf("/var/log/radiodan-%s.log", r.Name)
	return
}

func (r RadiodanApplication) PidFile() (path string) {
	path = fmt.Sprintf("/var/run/radiodan/radiodan-%s.pid", r.Name)
	return
}
