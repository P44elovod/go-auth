package handler

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	authHeader = "Authorization"
)

func (h *Handler) authMiddleWare(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authHeader)
		if header == "" {
			h.respondWithError(w, http.StatusUnauthorized, "empty auth header")

			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			h.respondWithError(w, http.StatusUnauthorized, "invalid auth header")

			return
		}

		if len(headerParts[1]) == 0 {
			h.respondWithError(w, http.StatusUnauthorized, "token is empty")

			return
		}

		userID, err := h.services.ParseAuthToken(headerParts[1])
		if err != nil {
			h.respondWithError(w, http.StatusUnauthorized, err.Error())

			return
		}
		fmt.Println(userID)
		next.ServeHTTP(w, r)
	})
}
