package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", checkIvanMessage)
	log.Println("Server started on: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func checkIvanMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command("python3", "train.py", string(body))
		stdout, err := cmd.Output()
		fmt.Fprintln(w, string(stdout))
	} else if r.Method == "GET" {
		fmt.Fprintf(w, "Hello world")
	}
}
