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

func (cfg *apiConfig) createFeedFollowHandle(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId string `json:"feed_id"`
	}

	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	feedFollowId := uuid.New()
	feedId, err := uuid.Parse(strings.TrimSpace(params.FeedId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}


	newFollowParams := database.CreateFeedFollowParams {
		ID: feedFollowId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID: feedId,
		UserID: user.ID,
	}
	newFeed, err := cfg.DB.CreateFeedFollow(r.Context(), newFollowParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, newFeed)
}

func (cfg *apiConfig) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}

func (cfg *apiConfig) delFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowId, err := uuid.Parse(strings.TrimSpace(r.PathValue("feedFollowId")))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	deleteParams := database.DeleteFeedFollowParams{
		ID: feedFollowId,
		UserID: user.ID,
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), deleteParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusAccepted, "")
}