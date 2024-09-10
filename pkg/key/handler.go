package key

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/J4yTr1n1ty/keyserver/pkg/internal/htmx"
	"github.com/J4yTr1n1ty/keyserver/pkg/internal/storage"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetKey(w http.ResponseWriter, r *http.Request) {
	emailParam := r.PathValue("email")
	key, err := storage.GetKey(emailParam)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	// Download public key
	w.Header().Set("Content-Disposition", "attachment; filename="+emailParam+".asc")
	w.Header().Set("Content-Type", "application/pgp-keys")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(key.PublicKey)))

	_, err = w.Write([]byte(key.PublicKey))
	if err != nil {
		http.Error(w, "Error writing key", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storage.ListAllKeys())
}

func (h *Handler) ListIdentities(w http.ResponseWriter, r *http.Request) {
	identityList := storage.GetUniqueIdentities()
	htmx.RenderListIdentities(w, identityList)
}

func (h *Handler) SubmitKey(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		htmx.RenderError(w, http.StatusBadRequest, "Error parsing form data")
		return
	}

	key := r.FormValue("key")
	if key == "" {
		htmx.RenderError(w, http.StatusBadRequest, "Key is required")
		return
	}

	_, err = storage.VerifyKey(key)
	if err != nil {
		htmx.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = storage.SubmitKey(key)
	if err != nil {
		htmx.RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	htmx.RenderSuccess(w, "Key submitted successfully")
}

func (h *Handler) VerifyMessage(w http.ResponseWriter, r *http.Request) {
	htmx.RenderError(w, http.StatusNotImplemented, "Not implemented yet")
}
