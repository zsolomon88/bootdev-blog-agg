package main

import "net/http"

func readinessHandle(w http.ResponseWriter, r *http.Request) {
	type readinessResponse struct {
		Status string `json:"status"`
	}

	respondWithJSON(w, 200, readinessResponse{Status: "ok"})
}


func errorHandle(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, "Internal Server Error")
}