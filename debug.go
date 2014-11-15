package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type DebugInfo struct {
	Timestamp          time.Time
	Hostname           string
	Addresses          []string
	InternetConnection bool
	Applications       []DebugApplication
}

type DebugApplication struct {
	Name string
	// TODO: this should be it's own struct
	DeployFile string
	LogFile    string
	PidFile    int64
	IsRunning  bool
}

func main() {
	port := parseFlags()
	http.HandleFunc("/", debugResponse)
	// TODO
	//http.HandleFunc("/download", downloadResponse)

	log.Printf("Debug server running on http://127.0.0.1:%d", port)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}

func debugResponse(w http.ResponseWriter, req *http.Request) {
	d := fetchDebugInfo()

	jsonReponse, err := json.Marshal(d)
	failOnError(err, "Cannot marshal JSON")

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonReponse)
}

func fetchDebugInfo() DebugInfo {
	d := DebugInfo{Timestamp: time.Now()}

	// fetch ip address(es)
	d.Hostname = hostname()
	d.Addresses = ipAddresses()
	// test for internet connection
	d.InternetConnection = checkConnection()
	// log active apps
	d.Applications = checkApps()

	return d
}

func hostname() (hostname string) {
	hostname, err := os.Hostname()
	failOnError(err, "Could not determine hostname")

	return
}

func ipAddresses() (ips []string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("[!] %s", err)
		return
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("[!] %s", err)
			return
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, ip.String())
		}
	}

	return
}

func checkConnection() bool {
	// match the request to a known output
	success := "<HTML><HEAD><TITLE>Success</TITLE></HEAD>" +
		"<BODY>Success</BODY></HTML>"

	res, err := http.Get("http://www.apple.com/library/test/success.html")

	if err != nil {
		return false
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return false
	}

	return string(body) == success
}

func checkApps() (apps []DebugApplication) {
	appNames := []string{
		"buttons", "cease", "example", "magic",
		"server", "updater", "debug",
	}

	for _, appName := range appNames {
		apps = append(apps, checkApp(appName))
	}

	return
}

func checkApp(appName string) (app DebugApplication) {
	app.Name = appName
	app.DeployFile = fetchDeployFile(appName)
	app.LogFile = fetchLogFile(appName)
	app.PidFile = fetchPidFile(appName)
	app.IsRunning = checkProcess(app.PidFile)
	return
}

func fetchDeployFile(appName string) (output string) {
	// TODO: pathTemplate should be a method on the struct
	pathTemplate := "/opt/radiodan/apps/%s/current/.deploy"
	path := fmt.Sprintf(pathTemplate, appName)

	file, err := ioutil.ReadFile(path)

	if err != nil {
		log.Println("[!] Could not open file", path)
		return
	}

	output = string(file)

	return
}

func fetchLogFile(appName string) (output string) {
	// TODO: logpath should be a method on the struct
	path := fmt.Sprintf("/var/log/radiodan-%s.log", appName)
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

func fetchPidFile(appName string) (output int64) {
	// TODO: pidpath should be a method on the struct
	pathTemplate := "/var/run/radiodan/radiodan-%s.pid"
	path := fmt.Sprintf(pathTemplate, appName)
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

func checkProcess(pid int64) (running bool) {
	running = false

	// this isn't the pid we're looking for
	if pid == 0 {
		return
	}

	// see if that process exists
	process, err := os.FindProcess(int(pid))
	if err != nil {
		return
	}

	// see if the process responds to sig 0
	err = process.Signal(syscall.Signal(0))
	running = (err == nil)

	return
}

func parseFlags() (port int) {
	flag.IntVar(&port, "port", 8080, "Port for server")
	flag.Parse()
	return
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
