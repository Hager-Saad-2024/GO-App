package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Answer struct {
	ID       string `json:"id,omitempty"`
	Answer_1 string `json:"answer1,omitempty"`
	Answer_2 string `json:"answer2,omitempty"`
	Answer_3 string `json:"answer3,omitempty"`
}

var client *mongo.Client

func main() {
	// Get server port from environment variable or use default
	serverPort := getEnv("SERVER_PORT", "8080")

	// Get MongoDB URI from environment variable or use default
	mongoURI := getEnv("MONGO_URI", "mongodb://mongo-local:27017")

	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Println("MongoDB not reachable at startup, readiness probe will handle it")
	}

	log.Println("Connected to MongoDB successfully")

	router := mux.NewRouter()

	// Application routes
	router.HandleFunc("/", GetQuestion).Methods("GET")
	router.HandleFunc("/", SubmitAnswer).Methods("POST")

	// Health endpoints
	router.HandleFunc("/health", HealthHandler).Methods("GET")
	router.HandleFunc("/ready", ReadyHandler).Methods("GET")

	log.Printf("Starting server on port %s...\n", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}

// =======================
// Utility
// =======================

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// =======================
// Liveness Probe
// =======================

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// =======================
// Readiness Probe
// =======================

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Println("Readiness check failed: MongoDB not reachable")
		http.Error(w, "Database not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}

// =======================
// Application Handlers
// =======================

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("What is your favorite programming language framework?")
}

func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	var answer Answer
	_ = json.NewDecoder(r.Body).Decode(&answer)

	log.Printf("answer.Answer_1: %v", answer.Answer_1)
	log.Printf("answer.Answer_2: %v", answer.Answer_2)
	log.Printf("answer.Answer_3: %v", answer.Answer_3)

	collection := client.Database("surveyDB").Collection("answers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, bson.M{
		"Answer1": answer.Answer_1,
		"Answer2": answer.Answer_2,
		"Answer3": answer.Answer_3,
	})

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(result.InsertedID)
}