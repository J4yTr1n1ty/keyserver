package homepage

import (
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Homepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
