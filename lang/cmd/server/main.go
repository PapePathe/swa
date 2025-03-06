package main

import (
	"net/http"
	"swahili/lang/server"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", server.LiveInterpreter{})
	http.ListenAndServe(":2108", mux)
}
