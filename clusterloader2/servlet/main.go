package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	testFile := os.Getenv("TESTRESULTS")
	if testFile == "" {
		testFile = "/junit.xml"
	}
	data, err := ioutil.ReadFile(testFile)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("results not here"))
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(data)
}
