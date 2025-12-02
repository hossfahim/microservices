package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"rides/internal/database"
	"rides/internal/server"
)

func main() {
	mongoURI := getEnv("MONGO_URI", "mongodb://user_app:strong_app_password@rides-service-database:27017/user_db?authSource=ridenow_rides")
	db, err := database.InitMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	usersServiceURL := getEnv("USERS_SERVICE_URL", "http://localhost:3000")
	port := fmt.Sprintf(":%s", getEnv("PORT", "8080"))

	s := server.NewServer(db, usersServiceURL)

	log.Printf("🚀 Service Rides démarré sur le port %s", port)
	if err := http.ListenAndServe(port, s); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
