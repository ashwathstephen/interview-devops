package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"devops-interview/database"
)

const (
	defaultHTTPAddr  = ":8080"
	defaultPGConnStr = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	defaultRedisAddr = "localhost:6379"
	healthPath       = "/health"
	readinessTimeout = 3 * time.Second
)

type healthResponse struct {
	Status   string            `json:"status"`
	Postgres string            `json:"postgres,omitempty"`
	Redis    string            `json:"redis,omitempty"`
	Details  map[string]string `json:"details,omitempty"`
}

func main() {
	httpAddr := getEnv("HTTP_ADDR", defaultHTTPAddr)
	pgConnStr := getEnv("DATABASE_URL", defaultPGConnStr)
	redisAddr := getEnv("REDIS_ADDR", defaultRedisAddr)

	ctx := context.Background()

	pg := database.NewPostgresOrNil(ctx, pgConnStr)
	if pg != nil {
		defer pg.Close()
	}

	rdb := database.NewRedisWithPing(ctx, redisAddr)
	if rdb != nil {
		defer rdb.Close()
	}

	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := healthResponse{Status: "ok", Details: make(map[string]string)}

		checkCtx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
		defer cancel()

		if pg != nil {
			if err := pg.Ping(checkCtx); err != nil {
				resp.Status = "degraded"
				resp.Postgres = "unhealthy"
				resp.Details["postgres"] = err.Error()
			} else {
				resp.Postgres = "healthy"
			}
		} else {
			resp.Postgres = "not_configured"
		}

		if rdb != nil {
			if err := rdb.Ping(checkCtx); err != nil {
				if resp.Status == "ok" {
					resp.Status = "degraded"
				}
				resp.Redis = "unhealthy"
				resp.Details["redis"] = err.Error()
			} else {
				resp.Redis = "healthy"
			}
		} else {
			resp.Redis = "not_configured"
		}

		code := http.StatusOK
		if resp.Status == "degraded" {
			code = http.StatusServiceUnavailable
		}
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(resp)
	}

	http.HandleFunc(healthPath, healthHandler)
	log.Printf("listening on %s", httpAddr)
	if err := http.ListenAndServe(httpAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
