package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"usage-lakehouse/internal/handler"
	"usage-lakehouse/internal/model"
	"usage-lakehouse/internal/repository"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type JSONRequest struct {
	Data []model.UsageData `json:"data"`
}

type Response struct {
	Message  string `json:"message"`
	S3Object string `json:"s3_object,omitempty"`
}

var (
	secretCache, _ = secretcache.New()
)

type HealthResponse struct {
	Status string `json:"status"`
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := HealthResponse{
		Status: "Success",
	}
	json.NewEncoder(w).Encode(response)
}

func getDSN() (string, error) {
	userName := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	secretArn := os.Getenv("DB_MASTER_SECRET_ARN")
	if secretArn != "" {
		log.Printf("getting password from secret: %v\n", secretArn)
		result, err := secretCache.GetSecretString(secretArn)
		if err != nil {
			return "", fmt.Errorf("error retreiving database secret: %v", err)
		}
		var passwordData struct {
			Password string `json:"password"`
			Username string `json:"username"`
		}
		if err := json.Unmarshal([]byte(result), &passwordData); err != nil {
			return "", fmt.Errorf("error retrieving database secret: %v", err)
		}
		password = passwordData.Password
		if passwordData.Username != "" {
			userName = passwordData.Username
		}

	}
	host := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=allow pool_max_conns=10", userName, password, host, dbPort, database)
	return dsn, nil
}

func main() {
	dsn, err := getDSN()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbpool.Close()
	accountRepo := repository.NewAccountRepository(dbpool)
	accountHandler := handler.NewAccountHandler(accountRepo)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/healthz", healthCheck)
	r.Post("/accounts", accountHandler.CreateAccount)
	r.Get("/accounts/{id}", accountHandler.GetAccount)
	r.Put("/accounts/{id}", accountHandler.UpdateAccount)
	r.Delete("/accounts/{id}", accountHandler.DeleteAccount)
	r.Get("/accounts", accountHandler.ListAccounts)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
