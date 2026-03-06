package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var req tripReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.UserID == "" {
		http.Error(w, "Missing userID", http.StatusBadRequest)
		return
	}

	// Why we need to create a new client for each connection:
	// because if a service is down, we don't want to block the whole application
	// so we create a new client for each connection
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	tripPreview, err := tripService.Client.TripPreview(r.Context(), req.toProto())
	if err != nil {
		log.Printf("Error calling TripPreview: %v\n", err)
		http.Error(w, "Failed to get trip preview", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{
		Data: tripPreview,
	}

	writeJSON(w, http.StatusCreated, response)
}

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	var req startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.UserID == "" || req.RideFareID == "" {
		http.Error(w, "Missing userID or rideFareID", http.StatusBadRequest)
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	tripStart, err := tripService.Client.CreateTrip(r.Context(), req.toProto())
	if err != nil {
		log.Printf("Error calling CreateTrip: %v\n", err)
		http.Error(w, "Failed to start trip", http.StatusInternalServerError)
		return
	}

	

	response := contracts.APIResponse{
		Data: tripStart,
	}
	writeJSON(w, http.StatusCreated, response)
}
