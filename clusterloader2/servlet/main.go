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

func main() {

	http.HandleFunc("/results", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}
	results := make(map[string]interface{})
	testDir := os.Getenv("TESTRESULTS")
	if testDir == "" {
		testDir = "/"
	}

	err := filepath.Walk(testDir, func(path string, info fs.FileInfo, err error) error {
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

	responseBody, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("could not fetch results: %s", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}
