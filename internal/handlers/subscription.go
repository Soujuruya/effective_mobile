package handlers

import (
	"context"
	"effective_mobile/internal/models"
	"effective_mobile/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

// SubscriptionHandler handles subscription-related endpoints
type SubscriptionHandler struct {
	Service *service.SubscriptionService
}

// Create godoc
// @Summary Create a subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription data"
// @Success 201 {object} models.Subscription
// @Failure 400 {string} string "invalid JSON or invalid fields"
// @Failure 500 {string} string "internal server error"
// @Router /subscriptions/create [post]
func (h *SubscriptionHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if len(sub.StartDate) != 7 || len(sub.EndDate) != 7 {
		http.Error(w, "invalid date format, must be YYYY-MM", http.StatusBadRequest)
		return
	}

	if err := h.Service.Insert(ctx, &sub); err != nil {
		log.Printf("failed to insert subscription: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

// List godoc
// @Summary List subscriptions
// @Description List all subscriptions with pagination
// @Tags subscriptions
// @Produce json
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} models.Subscription
// @Failure 400 {string} string "invalid parameters"
// @Router /subscriptions/list [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := r.URL.Query()

	offset := 0
	limit := 10

	if offsetStr := params.Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		} else {
			http.Error(w, "invalid offset value", http.StatusBadRequest)
			return
		}
	}

	if limitStr := params.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		} else {
			http.Error(w, "invalid limit value", http.StatusBadRequest)
			return
		}
	}

	subs := h.Service.Select(ctx, limit, offset)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subs)
}

// Get godoc
// @Summary Get a subscription
// @Description Get subscription by user_id and service_name
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID (UUID)"
// @Param service_name query string true "Service Name"
// @Success 200 {object} models.Subscription
// @Failure 400 {string} string "missing or invalid parameters"
// @Failure 500 {string} string "internal server error"
// @Router /subscriptions/get [get]
func (h *SubscriptionHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	userIDStr := params.Get("user_id")
	name := params.Get("service_name")

	if userIDStr == "" || name == "" {
		http.Error(w, "missing user_id or service_name", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	sub, err := h.Service.SelectByNameAndUserID(ctx, name, userID)
	if err != nil {
		log.Printf("failed to get subscription: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		log.Printf("failed to encode subscription to JSON: %v", err)
	}
}

// Update godoc
// @Summary Update a subscription
// @Description Update subscription fields
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription data"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "invalid JSON or missing fields"
// @Failure 500 {string} string "internal server error"
// @Router /subscriptions/update [put]
func (h *SubscriptionHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if sub.Name == "" || sub.Price <= 0 || sub.UserID == uuid.Nil || sub.StartDate == "" {
		http.Error(w, "missing or invalid fields", http.StatusBadRequest)
		return
	}

	if err := h.Service.Update(ctx, sub); err != nil {
		log.Printf("failed to update subscription: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete godoc
// @Summary Delete a subscription
// @Description Delete subscription by user_id and service_name
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID (UUID)"
// @Param name query string true "Service Name"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "missing or invalid parameters"
// @Failure 500 {string} string "internal server error"
// @Router /subscriptions/delete [delete]
func (h *SubscriptionHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	userIDStr := params.Get("user_id")
	name := params.Get("name")

	if userIDStr == "" || name == "" {
		http.Error(w, "missing user_id or service_name", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.Service.Delete(ctx, name, userID); err != nil {
		log.Printf("failed to delete subscription: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SumPrice godoc
// @Summary Sum subscription prices
// @Description Calculate total subscription cost over a period, filtered by user_id and service_name
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID (UUID)"
// @Param service_name query string false "Service Name"
// @Param start_date query string false "Start date YYYY-MM"
// @Param end_date query string false "End date YYYY-MM"
// @Success 200 {object} map[string]int "Total price"
// @Failure 400 {string} string "invalid parameters"
// @Failure 500 {string} string "internal server error"
// @Router /subscriptions/sum [get]
func (h *SubscriptionHandler) SumPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	var userID uuid.UUID
	var err error

	userIDStr := params.Get("user_id")
	name := params.Get("service_name")
	startDate := params.Get("start_date")
	endDate := params.Get("end_date")

	if userIDStr != "" {
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
	}

	sum, err := h.Service.SumPrice(ctx, name, userID, startDate, endDate)
	if err != nil {
		log.Printf("failed to sum subscriptions: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]int{"total price": sum}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode JSON: %v", err)
	}
}
