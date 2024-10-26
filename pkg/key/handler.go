package key

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/ProtonMail/gopenpgp/v3/crypto"

	"github.com/J4yTr1n1ty/keyserver/pkg/internal"
	"github.com/J4yTr1n1ty/keyserver/pkg/internal/htmx"
	"github.com/J4yTr1n1ty/keyserver/pkg/internal/storage"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetKeylist(w http.ResponseWriter, r *http.Request) {
	accept_header := r.Header["Accept"]
	if slices.Contains(accept_header, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(storage.GetUniqueIdentities())
	} else {
		w.Header().Set("Content-Type", "text/html")
		htmx.RenderKeyTable(w, storage.GetUniqueIdentities())
	}
}

func (h *Handler) GetKey(w http.ResponseWriter, r *http.Request) {
	emailParam := r.PathValue("email")
	key, err := storage.GetKeyByEmail(emailParam)
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	pgpKey, err := crypto.NewKeyFromArmored(string(key.PublicKey))
	if err != nil {
		http.Error(w, "Error parsing key", http.StatusInternalServerError)
		return
	}

	armoredPublicKey, err := pgpKey.Armor()
	if err != nil {
		http.Error(w, "Error getting armored public key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+emailParam+".asc")
	w.Header().Set("Content-Type", "application/pgp-keys")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(armoredPublicKey)))

	_, err = w.Write([]byte(armoredPublicKey))
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
	accept_header := r.Header["Accept"]
	if slices.Contains(accept_header, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(identityList)
	} else {
		w.Header().Set("Content-Type", "text/html")
		htmx.RenderListIdentities(w, identityList)
	}
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
	err := r.ParseForm()
	if err != nil {
		htmx.RenderError(w, http.StatusBadRequest, "Error parsing form data")
		return
	}

	message := r.PostFormValue("message")
	if message == "" {
		htmx.RenderError(w, http.StatusBadRequest, "Message is required")
		return
	}

	signer_fingerprint := r.PostFormValue("signer")
	if signer_fingerprint == "" {
		htmx.RenderError(w, http.StatusBadRequest, "Signer is required")
		return
	}

	signatureCreationDate, err := internal.VerifyMessage(signer_fingerprint, message)
	if err != nil {
		htmx.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	htmx.RenderSuccess(w, "Message verified successfully (created: "+signatureCreationDate+")")
}
