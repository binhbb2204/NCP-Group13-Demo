//go:build !mysql

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"DB-Presentation/database/sqlite"
	"DB-Presentation/db"
	"DB-Presentation/handlers"
	mongopkg "DB-Presentation/mongo"
	"DB-Presentation/utils"
	"DB-Presentation/ws"

	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func main() {
	dbPath := "data/chat.db"

	d, err := db.OpenDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	fmt.Println("✅ Successfully connected to SQLite!")

	if db.NeedsMigration(d) {
		fmt.Println("⚠️  Database needs migration! Run: go run migrate.go")
		fmt.Println("   Continuing anyway...")
	}

	// Seed initial data (admin user)
	if err := sqlite.SeedData(d); err != nil {
		log.Printf("⚠️  Warning: Could not seed data: %v\n", err)
	}

	// load .env if present (simple parser)
	loadEnvFile(".env")

	// connect to Mongo if URI provided
	var mongoClientPtr *mongodriver.Client
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		mc, err := mongopkg.Connect(uri)
		if err != nil {
			log.Println("warning: could not connect to mongo:", err)
		} else {
			mongoClientPtr = mc
			fmt.Println("✅ Connected to MongoDB")
		}
	}

	router := mux.NewRouter()

	// Register handlers and WebSocket route (pass mongo client if available)
	handlers.RegisterRoutes(router, d, mongoClientPtr)
	router.HandleFunc("/ws/{userId}", ws.HandleWebSocket)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	fmt.Println("🚀 Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", utils.EnableCORS(router)))
}

// loadEnvFile loads simple KEY=VALUE pairs from a file into environment variables.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)
		os.Setenv(key, val)
	}
}
