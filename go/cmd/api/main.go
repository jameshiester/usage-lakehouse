package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"usage-lakehouse/internal/handler"
	"usage-lakehouse/internal/model"
	"usage-lakehouse/internal/repository"

	"github.com/aws/aws-secretsmanager-caching-go/secretcache"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/xitongsys/parquet-go-source/s3"
	"github.com/xitongsys/parquet-go/writer"
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

func parseCSV(r io.Reader) ([]model.UsageData, error) {
	reader := csv.NewReader(r)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Validate header
	if len(header) != 2 || header[0] != "asset_id" || header[1] != "usage_qty" {
		return nil, errors.New("invalid CSV format. expected columns: asset_id, usage_qty")
	}

	var data []model.UsageData
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		usageQty, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, errors.New("invalid usage_qty format")
		}

		data = append(data, model.UsageData{
			AssetID:  record[0],
			UsageQty: usageQty,
		})
	}

	return data, nil
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	accountId := chi.URLParam(r, "account_id")
	if accountId == "" {
		http.Error(w, "accountId is required in the path", http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	var data []model.UsageData
	var err error

	if format == "csv" {
		data, err = parseCSV(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// Default to JSON format
		var jsonReq JSONRequest
		if err := json.NewDecoder(r.Body).Decode(&jsonReq); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		data = jsonReq.Data
	}
	// Upload to S3
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		http.Error(w, "S3_BUCKET environment variable not set", http.StatusInternalServerError)
		return
	}

	key := accountId + "/usage_data_" + time.Now().Format("20060102_150405") + ".parquet"
	fw, err := s3.NewS3FileWriter(r.Context(), bucket, key, "bucket-owner-full-control", nil)
	if err != nil {
		log.Println("Can't open file", err)
		return
	}
	pw, err := writer.NewParquetWriter(fw, new(model.UsageData), 4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}
	for _, usage := range data {
		if err = pw.Write(usage); err != nil {
			log.Println("Write error", err)
		}
	}
	if err = pw.WriteStop(); err != nil {
		log.Println("WriteStop err", err)
	}
	err = fw.Close()
	if err != nil {
		log.Println("Error closing S3 file writer")
	}

	response := Response{
		Message:  "Data written to S3 successfully",
		S3Object: key,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

func main() {
	userName := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	secretArn := os.Getenv("DB_MASTER_SECRET_ARN")
	if secretArn != "" {
		fmt.Printf("getting password from secret: %v\n", secretArn)
		result, err := secretCache.GetSecretString(secretArn)
		if err != nil {
			fmt.Printf("error retrieving database secret: %v\n", err)
			log.Fatal("error retrieving database secret")
		}
		var passwordData struct {
			Password string `json:"password"`
			Username string `json:"username"`
		}
		if err := json.Unmarshal([]byte(result), &passwordData); err != nil {
			fmt.Printf("error retrieving database secret: %v\n", err)
			log.Fatal("error unmarshaling password")
		}
		password = passwordData.Password
		if passwordData.Username != "" {
			userName = passwordData.Username
		}

	}
	host := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=verify-ca pool_max_conns=10", userName, password, host, dbPort, database)
	// dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", userName, password, host, dbPort, database)
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer dbpool.Close()
	accountRepo := repository.NewAccountRepository(dbpool)
	accountHandler := handler.NewAccountHandler(accountRepo)

	r := chi.NewRouter()
	r.Get("/healthz", healthCheck)
	r.Post("/accounts", accountHandler.CreateAccount)
	r.Get("/accounts/{id}", accountHandler.GetAccount)
	r.Put("/accounts/{id}", accountHandler.UpdateAccount)
	r.Delete("/accounts/{id}", accountHandler.DeleteAccount)
	r.Get("/accounts", accountHandler.ListAccounts)
	r.Post("/edi/monthly-usage", helloWorldHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
