package server

import (
	"net/http"

	"github.com/topritchett/game-server/proxmox"
)

// Handler for http requests
type Handler struct {
	mux *http.ServeMux
}

// New http handler
func New(s *http.ServeMux) *Handler {
	h := Handler{s}
	h.registerRoutes()

	return &h
}

// RegisterRoutes for all http endpoints
func (h *Handler) registerRoutes() {
	h.mux.HandleFunc("/", h.HelloWorld)
	h.mux.HandleFunc("/proxurl", h.ServerGetProxUrl)
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Hello World"))
}

func (h *Handler) ServerGetProxUrl(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Write([]byte(proxmox.GetProxUrl(proxmox.Auth, proxmox.QemuUrl)))
}
