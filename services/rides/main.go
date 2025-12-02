package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Ride struct {
	ID            string    `json:"id"`
	PassengerID   string    `json:"passengerId"`
	DriverID      string    `json:"driverId"`
	FromZone      string    `json:"from_zone"`
	ToZone        string    `json:"to_zone"`
	Price         float64   `json:"price"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"paymentStatus"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

var rides = make(map[string]*Ride)
var mu sync.Mutex

func main() {
	http.HandleFunc("/rides", createRide)
	http.HandleFunc("/rides/", rideHandler)

	fmt.Println("Ride service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createRide(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PassengerID string `json:"passengerId"`
		FromZone    string `json:"from_zone"`
		ToZone      string `json:"to_zone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	price := getPrice(req.FromZone, req.ToZone)
	driverID := getAvailableDriver()

	rideID := uuid.New().String()
	ride := &Ride{
		ID:            rideID,
		PassengerID:   req.PassengerID,
		DriverID:      driverID,
		FromZone:      req.FromZone,
		ToZone:        req.ToZone,
		Price:         price,
		Status:        "ASSIGNED",
		PaymentStatus: "PENDING",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	mu.Lock()
	rides[rideID] = ride
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ride)
}

func rideHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/rides/")
	parts := strings.Split(path, "/")

	id := parts[0]

	mu.Lock()
	ride, ok := rides[id]
	mu.Unlock()

	if !ok {
		http.Error(w, "Ride not found", http.StatusNotFound)
		return
	}

	// GET /rides/{id}
	if len(parts) == 1 && r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ride)
		return
	}

	// PATCH /rides/{id}/status
	if len(parts) == 2 && parts[1] == "status" && r.Method == http.MethodPatch {

		var req struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ride.Status = req.Status
		ride.UpdatedAt = time.Now()

		if req.Status == "COMPLETED" {
			ride.PaymentStatus = capturePayment(ride.ID, ride.Price)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ride)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func getPrice(from, to string) float64 {
	return float64(rand.Intn(50) + 10)
}

func getAvailableDriver() string {
	return fmt.Sprintf("driver-%d", rand.Intn(100))
}

func capturePayment(rideID string, amount float64) string {
	log.Printf("Payment captured for ride %s: %.2f\n", rideID, amount)
	return "CAPTURED"
}
