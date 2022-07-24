package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gsockets/gsockets"
	appmanagers "github.com/gsockets/gsockets/app_managers"
)

type AuthMiddleware struct {
	apps gsockets.AppManager
}

func NewAuthMiddleware(apps gsockets.AppManager) *AuthMiddleware {
	return &AuthMiddleware{apps: apps}
}

// Handler verifies the existence of the app and the signature for the request.
// See https://pusher.com/docs/channels/library_auth_reference/rest-api#Authentication for implementation details.
func (auth *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app, err := auth.apps.FindById(r.Context(), chi.URLParam(r, "appId"))
		if err != nil {
			if errors.Is(err, appmanagers.ErrInvalidAppId) {
				RenderJSON(w, http.StatusBadRequest, err.Error(), nil)
				return
			}

			RenderJSON(w, http.StatusInternalServerError, "internal server error", nil)
			return 
		}

		queryParams := r.URL.Query()
		keys := make([]string, 0)
		authSignature := queryParams.Get("auth_signature")

		for key := range queryParams {
			if key == "auth_signature" {
				continue
			}

			keys = append(keys, key)
		}

		sort.Strings(keys)

		queryParamsSorted := make([]string, 0)
		var signatureString strings.Builder

		signatureString.WriteString(r.Method)
		signatureString.WriteString("\n")
		signatureString.WriteString(r.URL.Path)
		signatureString.WriteString("\n")

		for _, key := range keys {
			var str strings.Builder
			str.WriteString(key)
			str.WriteString("=")
			str.WriteString(queryParams.Get(key))

			queryParamsSorted = append(queryParamsSorted, str.String())
		}

		signatureString.WriteString(strings.Join(queryParamsSorted, "&"))

		incomingSignature, err := hex.DecodeString(authSignature)
		if err != nil {
			RenderJSON(w, http.StatusUnauthorized, "invalid signature string", nil)
			return
		}

		hasher := hmac.New(sha256.New, []byte(app.Secret))
		hasher.Write([]byte(signatureString.String()))

		if valid := hmac.Equal(hasher.Sum(nil), incomingSignature); !valid {
			RenderJSON(w, http.StatusUnauthorized, "invalid signature string", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
