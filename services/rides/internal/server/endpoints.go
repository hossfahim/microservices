package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"rides/internal/types"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) createRide(w http.ResponseWriter, r *http.Request) {
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
	driverID, err := s.getAvailableDriver()
	if err != nil {
		log.Printf("[ERROR] Failed to get available driver: %v", err)
		http.Error(w, "No available driver found", http.StatusServiceUnavailable)
		return
	}

	ride := &types.Ride{
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.db.CreateRide(ctx, ride)
	if err != nil {
		log.Printf("[ERROR] Failed to create ride: %v", err)
		http.Error(w, "Error creating ride", http.StatusInternalServerError)
		return
	}

	if err := s.updateDriverStatus(driverID, false); err != nil {
		log.Printf("[WARN] Failed to update driver status: %v", err)
	}

	log.Printf("[CREATE] Nouvelle course créée: ID=%s, Passenger=%s, Driver=%s", ride.ID.Hex(), ride.PassengerID, ride.DriverID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ride)
}

func (s *Server) getRide(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ride, err := s.db.GetRideByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Ride not found", http.StatusNotFound)
			return
		}
		log.Printf("[ERROR] Failed to get ride: %v", err)
		http.Error(w, "Error retrieving ride", http.StatusInternalServerError)
		return
	}

	log.Printf("[READ] Course récupérée: ID=%s", idStr)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ride)
}

func (s *Server) updateRideStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.db.UpdateRideStatus(ctx, id, req.Status)
	if err != nil {
		log.Printf("[ERROR] Failed to update ride status: %v", err)
		http.Error(w, "Error updating ride status", http.StatusInternalServerError)
		return
	}

	// If status is COMPLETED, capture payment
	if req.Status == "COMPLETED" {
		ride, err := s.db.GetRideByID(ctx, id)
		if err == nil {
			paymentStatus := capturePayment(ride.ID.Hex(), ride.Price)
			err = s.db.UpdateRidePaymentStatus(ctx, id, paymentStatus)
			if err != nil {
				log.Printf("[ERROR] Failed to update payment status: %v", err)
			}
		}

		if err := s.updateDriverStatus(ride.DriverID, true); err != nil {
			log.Printf("[WARN] Failed to update driver status: %v", err)
		}
	}

	// Get updated ride to return
	ride, err := s.db.GetRideByID(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Failed to get updated ride: %v", err)
		http.Error(w, "Error retrieving updated ride", http.StatusInternalServerError)
		return
	}

	log.Printf("[UPDATE] Statut de la course %s mis à jour: %s", idStr, req.Status)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ride)
}

func getPrice(from, to string) float64 {
	return float64(rand.Intn(50) + 10)
}

func (s *Server) getAvailableDriver() (string, error) {
	url := fmt.Sprintf("%s/drivers?available=true", s.usersServiceURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call users service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("users service returned status %d: %s", resp.StatusCode, string(body))
	}

	var drivers []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		IsAvailable bool   `json:"is_available"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&drivers); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(drivers) == 0 {
		return "", fmt.Errorf("no available drivers found")
	}

	// Return the first available driver's ID
	return drivers[0].ID, nil
}

func (s *Server) updateDriverStatus(driverID string, isAvailable bool) error {
	url := fmt.Sprintf("%s/drivers/%s/status", s.usersServiceURL, driverID)

	payload := struct {
		IsAvailable bool `json:"is_available"`
	}{
		IsAvailable: isAvailable,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "PATCH", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call users service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("users service returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[UPDATE] Driver %s availability set to %v", driverID, isAvailable)
	return nil
}

func capturePayment(rideID string, amount float64) string {
	log.Printf("Payment captured for ride %s: %.2f\n", rideID, amount)
	return "CAPTURED"
}
