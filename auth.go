package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zsolomon88/bootdev-blog-agg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		apiKey = strings.TrimPrefix(apiKey, "ApiKey ")

		foundUser, err := cfg.DB.GetUserByKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Access Denied: %v", err))
			return 
		}

		handler(w, r, foundUser)
	}
}