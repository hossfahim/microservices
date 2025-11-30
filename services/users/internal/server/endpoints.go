package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"users/internal/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) createDriver(w http.ResponseWriter, r *http.Request) {
	var driver types.Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	driver.IsAvailable = true // Par défaut disponible
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.CreateDriver(ctx, &driver)
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}

	// Log clair pour l'observabilité
	log.Printf("[CREATE] Nouveau chauffeur créé: %s (ID: %s)", driver.Name, driver.ID.Hex())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(driver)
}

// getDrivers : Liste les chauffeurs (filtre optionnel ?available=true)
func (s *Server) getDrivers(w http.ResponseWriter, r *http.Request) {
	availableQuery := r.URL.Query().Get("available")

	var available *bool
	if availableQuery == "true" {
		val := true
		available = &val
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	drivers, err := s.db.GetDrivers(ctx, available)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log pour tracer les appels inter-services (Pricing ou Ride qui cherche un driver)
	log.Printf("[READ] Recherche drivers (available=%s) -> %d trouvés", availableQuery, len(drivers))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

// setStatus : Change la disponibilité (ex: quand une course est assignée)
func (s *Server) setStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // Go 1.22 feature
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "ID Invalide", http.StatusBadRequest)
		return
	}

	// Structure simple pour recevoir le status
	var statusUpdate struct {
		IsAvailable bool `json:"is_available"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.db.UpdateDriverStatus(ctx, id, statusUpdate.IsAvailable)
	if err != nil {
		http.Error(w, "Erreur Update", http.StatusInternalServerError)
		return
	}

	log.Printf("[UPDATE] Chauffeur %s -> disponibilité: %v", idStr, statusUpdate.IsAvailable)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) createPassenger(w http.ResponseWriter, r *http.Request) {
	var passenger types.Passenger
	if err := json.NewDecoder(r.Body).Decode(&passenger); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.CreatePassenger(ctx, &passenger)
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}

	log.Printf("[CREATE] Nouveau passager créé: %s (ID: %s)", passenger.Name, passenger.ID.Hex())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(passenger)
}

// getPassengers : Liste tous les passagers
func (s *Server) getPassengers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	passengers, err := s.db.GetPassengers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[READ] Recherche passagers -> %d trouvés", len(passengers))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passengers)
}

// getPassenger : Récupère un passager par son ID
func (s *Server) getPassenger(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "ID Invalide", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	passenger, err := s.db.GetPassengerByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Passager non trouvé", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passenger)
}

// updatePassenger : Met à jour un passager
func (s *Server) updatePassenger(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "ID Invalide", http.StatusBadRequest)
		return
	}

	var passenger types.Passenger
	if err := json.NewDecoder(r.Body).Decode(&passenger); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.db.UpdatePassenger(ctx, id, &passenger)
	if err != nil {
		http.Error(w, "Erreur Update", http.StatusInternalServerError)
		return
	}

	log.Printf("[UPDATE] Passager %s mis à jour", idStr)

	// Récupérer le passager mis à jour pour le retourner
	updatedPassenger, err := s.db.GetPassengerByID(ctx, id)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPassenger)
}

// deletePassenger : Supprime un passager
func (s *Server) deletePassenger(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "ID Invalide", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.db.DeletePassenger(ctx, id)
	if err != nil {
		http.Error(w, "Erreur Delete", http.StatusInternalServerError)
		return
	}

	log.Printf("[DELETE] Passager %s supprimé", idStr)
	w.WriteHeader(http.StatusNoContent)
}
