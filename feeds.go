package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (cfg *apiConfig) getSingleFeed(w http.ResponseWriter, r *http.Request, ) {
	feedId, err := uuid.Parse(strings.TrimSpace(r.PathValue("feedId")))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
		return
	}

	getParams := database.UpdateFetchTimeParams{
		ID: feedId,
		FetchedAt: sql.NullTime{
			Time: time.Now().UTC(),
			Valid: true,
		},
	}
	feed, err := cfg.DB.UpdateFetchTime(r.Context(), getParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, feed)
}


func (cfg *apiConfig) getNextFeeds(w http.ResponseWriter, r *http.Request, ) {
	numberOfFeeds := r.URL.Query().Get("limit")
	feedLimit := 10
	if numberOfFeeds != "" {
		feedLimit, _ = strconv.Atoi(numberOfFeeds)
	}

	feeds, err := cfg.DB.FetchNextFeeds(r.Context(), int32(feedLimit))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
		return
	}

	rssData := make([]Rss, feedLimit)
	type responseData struct {
		Titles []string `json:"titles"`
	}

	titleList := responseData{}
	for i, feed := range feeds {
		timeParam := database.UpdateFetchTimeParams{
			ID: feed.ID,
			FetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		}
		_, err := cfg.DB.UpdateFetchTime(r.Context(), timeParam)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
			return
		}

		feeds[i].FetchedAt = timeParam.FetchedAt
		var feedData *Rss = &Rss{}
		err = fetchRssFeed(feed.Url, feedData)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Internal Server Error: %v", err))
			return
		}
		rssData = append(rssData, *feedData)
		titleList.Titles = append(titleList.Titles, feedData.Channel.Title)
	}

	respondWithJSON(w, http.StatusOK, titleList)
}