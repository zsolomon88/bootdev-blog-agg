package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zsolomon88/bootdev-blog-agg/internal/database"
)

func (cfg *apiConfig) createUserHandle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	userName := strings.TrimSpace(params.Name)
	userId := uuid.New()

	ctx := r.Context()
/*
	randomData := make([]byte, 32)
	rand.Read(randomData)
	apiKey := hex.EncodeToString(randomData)
*/
	userParams := database.CreateUserParams{
		ID: userId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: userName,
	}
	insertedUser, err := cfg.DB.CreateUser(ctx, userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database Error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, insertedUser)

}

func (cfg *apiConfig) validateApiKeyHandle(w http.ResponseWriter, r *http.Request, user database.User) {

	respondWithJSON(w, http.StatusOK, user)
}