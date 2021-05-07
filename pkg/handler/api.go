package handler

import "net/http"

func (h *Handler) helloAPI(w http.ResponseWriter, r *http.Request) {
	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"palce": "api",
	})
}
