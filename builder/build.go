package builder

import (
	"encoding/json"
	"github.com/radiodan/debug/radiodan"
	"github.com/radiodan/debug/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var appNames = []string{
	"buttons", "cease", "example", "magic",
	"server", "updater", "debug",
}

func Build() radiodan.DebugInfo {
	d := radiodan.DebugInfo{Timestamp: time.Now()}

	// fetch ip address(es)
	d.Hostname = utils.Hostname()
	d.Addresses = utils.IpAddresses()
	// test for internet connection
	d.InternetConnection = utils.CheckConnection()
	// log active apps
	d.Applications = checkApps()

	return d
}

func checkApps() (apps []radiodan.RadiodanApplication) {
	for _, appName := range appNames {
		apps = append(apps, checkApp(appName))
	}

	return
}

func checkApp(appName string) (app radiodan.RadiodanApplication) {
	app.Name = appName
	app.Deploy = fetchDeployFile(app)
	app.LogTail = fetchLogFile(app)
	app.Pid = fetchPidFile(app)
	app.IsRunning = utils.CheckProcess(app.Pid)
	return
}

func fetchDeployFile(app radiodan.RadiodanApplication) (output radiodan.Deploy) {
	file, err := ioutil.ReadFile(app.DeployFile())

	if err != nil {
		log.Println("[!] Could not open file", app.DeployFile())
		return
	}

	err = json.Unmarshal(file, output)

	return
}

func fetchLogFile(app radiodan.RadiodanApplication) (output string) {
	path := app.LogFile()
	_, err := os.Stat(path)

	if err != nil {
		log.Printf("[!] Could not open file %s", path)
		return
	}

	stdout, err := exec.Command("/usr/bin/tail", "-n 100", path).Output()

	if err != nil {
		log.Printf("[!] Command failed: %s", err)
	} else {
		output = string(stdout)
	}

	return
}

func fetchPidFile(app radiodan.RadiodanApplication) (output int64) {
	path := app.PidFile()
	file, err := ioutil.ReadFile(path)

	if err != nil {
		log.Println("[!] Could not open file", path)
		return
	}

	pidString := strings.Trim(string(file), "\n")
	output, err = strconv.ParseInt(pidString, 10, 0)

	if err != nil {
		log.Println("[!] Could not parse pid as integer", pidString)
	}

	return
}
