package main

import (
	"fmt"
	"net/http"
)

func main() {
	run()
}

const welcomeMessage = "Hello Paack"

func run() {
	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, welcomeMessage)
	})

	http.ListenAndServe(":8080", nil)
}
