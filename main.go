package main

import (
	"net/http"

	"github.com/topritchett/game-server/server"
)

func main() {
	mux := http.NewServeMux()
	server.New(mux)

	http.ListenAndServe(":8080", mux)
}
