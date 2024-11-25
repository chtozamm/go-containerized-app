package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const (
	countKey      = "count"
	serverAddress = "0.0.0.0:3000"
)

func main() {
	REDIS_HOST := os.Getenv("REDIS_HOST")
	REDIS_PORT := os.Getenv("REDIS_PORT")
	if REDIS_HOST == "" || REDIS_PORT == "" {
		fmt.Println("Error: REDIS_HOST and REDIS_PORT must be set in environment variables")
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: REDIS_HOST + ":" + REDIS_PORT,
	})
	defer rdb.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Attempt to get the current count
		val, err := rdb.Get(ctx, countKey).Result()
		if err != nil {
			if err == redis.Nil {
				// If the key does not exist, set it to 0
				if err := rdb.Set(ctx, countKey, 0, 0).Err(); err != nil {
					http.Error(w, "Error setting count", http.StatusInternalServerError)
					return
				}
				val = "0"
			} else {
				// Handle other potential errors from Get
				http.Error(w, "Error getting count from db", http.StatusInternalServerError)
				return
			}
		}

		// Convert the value to an integer
		count, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, "Error converting count to integer", http.StatusInternalServerError)
			return
		}

		// Increment the count
		count++
		if err := rdb.Set(ctx, countKey, count, 0).Err(); err != nil {
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
	})

	fmt.Println("Server is listening on", serverAddress)
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
