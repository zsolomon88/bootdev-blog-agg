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

func (cfg *apiConfig) createFeedHandle(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}

	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	feedId := uuid.New()
	feedName := strings.TrimSpace(params.Name)
	feedUrl := strings.TrimSpace(params.Url)

	createParams := database.CreateFeedParams {
		ID: feedId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: feedName,
		Url: feedUrl,
		UserID: user.ID,
	}
	newFeed, err := cfg.DB.CreateFeed(r.Context(), createParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
		return
	}

	feedFollowId := uuid.New()
	newFollowParams := database.CreateFeedFollowParams {
		ID: feedFollowId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID: feedId,
		UserID: user.ID,
	}
	newFeedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), newFollowParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request: %v", err))
		return
	}
	type responseJson struct {
		Feed database.Feed `json:"feed"`
		FeedFollow database.Follow `json:"feed_follow"`
	}

	newFeedResponse := responseJson {
		Feed: newFeed,
		FeedFollow: newFeedFollow,
	}
	respondWithJSON(w, http.StatusCreated, newFeedResponse)
}

func (cfg *apiConfig) getAllFieldsHandle(w http.ResponseWriter, r *http.Request, ) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, feeds)
}