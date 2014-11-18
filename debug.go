package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/radiodan/debug/builder"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/debug", viewResponse)
	http.HandleFunc("/debug/download", downloadResponse)

	port := parseFlags()
	log.Printf("Debug server running on http://127.0.0.1:%d", port)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}

func viewResponse(w http.ResponseWriter, req *http.Request) {
	d := builder.Build()

	jsonResponse, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func downloadResponse(w http.ResponseWriter, req *http.Request) {
	d := builder.Build()
	jsonResponse, err := json.Marshal(d)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create zip bufer
	buf := new(bytes.Buffer)

	// Create zip archive.
	zw := zip.NewWriter(buf)
	f, err := zw.Create("debug.json")
	f.Write([]byte(jsonResponse))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, app := range d.Applications {
		filePath := app.LogFile()
		fileElem := strings.Split("/", filePath)
		fileName := fileElem[len(fileElem)-1]

		file, err := os.Open(filePath)
		defer file.Close()

		if err != nil {
			log.Println("[!] " + err.Error())
			continue
		}

		zipWriter, err := zw.Create(fileName)

		if err != nil {
			log.Println("[!] " + err.Error())
			continue
		}
		_, err = io.Copy(zipWriter, file)

		if err != nil {
			log.Println("[!] " + err.Error())
			continue
		}
	}

	err = zw.Close()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/x-gzip")
	w.Header().Set("Content-Disposition", "inline; filename=\"debug.zip\"")
	w.Write(buf.Bytes())
}

func parseFlags() (port int) {
	flag.IntVar(&port, "port", 8080, "Port for server")
	flag.Parse()
	return
}
