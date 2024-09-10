package api

import (
	"net/http"

	"github.com/J4yTr1n1ty/keyserver/pkg/homepage"
	"github.com/J4yTr1n1ty/keyserver/pkg/key"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	homepageHandler := homepage.NewHandler()

	mux.HandleFunc("GET /", homepageHandler.Homepage)

	keyHandler := key.NewHandler()

	mux.HandleFunc("GET /identities", keyHandler.ListIdentities)
	mux.HandleFunc("GET /list-all", keyHandler.ListAll)
	mux.HandleFunc("GET /key/{email}", keyHandler.GetKey)
	mux.HandleFunc("POST /submit-key", keyHandler.SubmitKey)
	mux.HandleFunc("POST /verify-message", keyHandler.VerifyMessage)

	return mux
}
