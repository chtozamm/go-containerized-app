package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const (
	countKey      = "count"
	serverAddress = ":3000"
)

type application struct {
	rdb *redis.Client
	ID  string
}

func main() {
	app := new(application)
	app.ID = os.Getenv("APP_ID")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" || redisPort == "" {
		fmt.Println("Error: REDIS_HOST and REDIS_PORT must be set in environment variables")
		os.Exit(1)
	}

	app.rdb = redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})
	defer app.rdb.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.rootHandler)

	srv := http.Server{
		Addr:              serverAddress,
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	fmt.Println("Server is listening on", serverAddress)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}

func (app *application) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Println("Request was served by app", app.ID)

	// Attempt to get the current count
	val, err := app.rdb.Get(ctx, countKey).Result()
	if err != nil {
		if err == redis.Nil {
			// If the key does not exist, set it to 0
			if err := app.rdb.Set(ctx, countKey, 0, 0).Err(); err != nil {
				http.Error(w, "Error setting count to 0", http.StatusInternalServerError)
				return
			}
			val = "0"
		} else {
			http.Error(w, "Error getting count from db", http.StatusInternalServerError)
			return
		}
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		http.Error(w, "Error converting count to integer", http.StatusInternalServerError)
		return
	}

	// Increment the count
	count++
	if err := app.rdb.Set(ctx, countKey, count, 0).Err(); err != nil {
		http.Error(w, "Error incrementing count in db", http.StatusInternalServerError)
		return
	}

	var times string
	if count == 1 {
		times = "time"
	} else {
		times = "times"
	}

	fmt.Fprintf(w, "Hello there! This page was requested %d %s.", count, times)
}
