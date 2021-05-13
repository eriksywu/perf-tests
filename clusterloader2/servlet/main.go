package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var status string

func main() {

	http.HandleFunc("/results", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	if status != "" {
		return
	}
	s, k := r.URL.Query()["status"]
	if !k || len(s) == 0 {
		return
	}
	fmt.Printf("status updated to ", s[0])
	status = s[0]
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		doneHandler(w, r)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	testDir := os.Getenv("TESTRESULTS")
	fmt.Printf("watching over testdir %s \n", testDir)
	if testDir == "" {
		testDir = "/"
	}
	if status == "" {
		w.WriteHeader(400)
		w.Write([]byte("not done"))
		return
	}
	results := make(map[string]interface{})
	results["status"] = status
	err := filepath.Walk(testDir, func(path string, info fs.FileInfo, err error) error {
		fmt.Printf("path = %s \n", path)
		if info == nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".json") {
			return nil
		}
		if strings.Contains(path, "metadata") {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		split := strings.Split(filepath.Base(path), "_")
		if len(split) < 2 {
			return nil
		}
		testName := split[0]
		var jsonData interface{}

		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			return err
		}
		results[testName] = jsonData
		return nil
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("could not fetch results: %s", err)))
		return
	}
	responseBody, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("could not fetch results: %s", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
