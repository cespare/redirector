package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	addr = flag.String("addr", "localhost:9310", "Listen addr")
	configFile = flag.String("conf", "sample.conf", "Config file")
)

type config map[string]string

func loadConfig(path string) (config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := config{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) != 2 {
			return nil, fmt.Errorf("Wrong number of fields (expected 2, got %d)", len(parts))
		}
		key, value := parts[0], parts[1]
		if _, ok := c[key]; ok {
			return nil, fmt.Errorf("Duplicate key in config: %s", key)
		}
		c[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return c, nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Load the file each time so we pick up changes.
	c, err := loadConfig(*configFile)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Redirector is misconfigured.", http.StatusInternalServerError)
		return
	}

	// Normalize path
	path := strings.Trim(r.URL.Path, "/")

	target, ok := c[path]
	if !ok {
		log.Printf("Bad request for '%s'", path)
		http.Error(w, "No such file.", http.StatusNotFound)
		return
	}

	log.Printf("Redirecting %s to %s.", path, target)
	http.Redirect(w, r, target, http.StatusFound)
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)
	log.Println("Listening on", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
