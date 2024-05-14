package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zsolomon88/bootdev-blog-agg/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
	DbUrl string
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("CONN")


	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	cfg := apiConfig {
		DB: dbQueries,
		DbUrl: dbUrl,
	}

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("GET /v1/readiness", readinessHandle)
	httpMux.HandleFunc("GET /v1/err", errorHandle)
	httpMux.HandleFunc("POST /v1/users", cfg.createUserHandle)
	httpMux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.validateApiKeyHandle))
	httpMux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.createFeedHandle))
	httpMux.HandleFunc("GET /v1/feeds", cfg.getAllFieldsHandle)
	
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%v", os.Getenv("PORT")),
		Handler: httpMux,
	}

	log.Printf("Starting server on port: %v\n", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())
}