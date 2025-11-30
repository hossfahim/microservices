package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"users/internal/database"
	"users/internal/server"
)

func main() {
	mongoURI := getEnv("MONGO_URI", "mongodb://user_app:strong_app_password@users-service-database:27017/user_db?authSource=ridenow_users")
	db, err := database.InitMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	port := fmt.Sprintf(":%s", getEnv("PORT", "3000"))

	s := server.NewServer(db)

	log.Printf("🚀 Service Users démarré sur le port %s", port)
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
