package http

import (
	"encoding/json"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/types"
)

type Handler struct {
	Service domain.TripService
}

type PreviewRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (h *Handler) HandlePreview(w http.ResponseWriter, r *http.Request) {
	var req PreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	trip, err := h.Service.GetRoute(ctx, &req.Pickup, &req.Destination)
	if err != nil {
		http.Error(w, "failed to get route", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{
		Data: trip,
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
