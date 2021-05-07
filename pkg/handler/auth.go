package handler

import (
	"encoding/json"
	"net/http"

	"github.com/p44elovod/auth-with-gopg/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Errorf("error occured at decode body: %s", err)
		h.respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}

	id, err := h.services.CreateUser(input)
	if err != nil {
		logrus.Errorf("error occured at user creation stage: %s", err)
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required" pg:"password_hash"`
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Errorf("error occured at decode body: %s", err)
		h.respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}

	tokens, err := h.services.GenerateTokenPair(input.Username, input.Password)
	if err != nil {
		logrus.Errorf("error occured at auth token pair generation: %s", err)
		h.respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"authToken":    tokens["authToken"],
		"refreshToken": tokens["refreshToken"],
	})
}

type refreshInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (h *Handler) refreshTokenPair(w http.ResponseWriter, r *http.Request) {
	var input refreshInput
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		logrus.Errorf("error occured at decode body: %s", err)
		h.respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}

	tokens, err := h.services.RefreshTokens(input.RefreshToken)
	if err != nil {
		logrus.Errorf("error occured at auth token pair refreshing: %s", err)
		h.respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"authToken":    tokens["authToken"],
		"refreshToken": tokens["refreshToken"],
	})
}
