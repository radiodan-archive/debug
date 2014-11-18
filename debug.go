package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/radiodan/debug/debug"
	"log"
	"net/http"
)

func main() {
	port := parseFlags()
	http.HandleFunc("/", debugResponse)
	//http.HandleFunc("/download", downloadResponse)

	log.Printf("Debug server running on http://127.0.0.1:%d", port)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}

func debugResponse(w http.ResponseWriter, req *http.Request) {
	d := debug.Build()

	jsonReponse, err := json.Marshal(d)
	if err != nil {
		msg := "Cannot marshal JSON"
		log.Fatalf("%s: %s", msg, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonReponse)
}

func parseFlags() (port int) {
	flag.IntVar(&port, "port", 8080, "Port for server")
	flag.Parse()
	return
}
