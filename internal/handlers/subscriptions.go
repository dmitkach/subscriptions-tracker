package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Subscription struct {
	ID          uuid.UUID `json:"id,omitempty"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type Handler struct {
	DB *sql.DB
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var sub Subscription
	if err := json.Unmarshal(body, &sub); err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		http.Error(w, "incorrect start_date format, expected: MM-YYYY", http.StatusBadRequest)
		return
	}

	var endDate *time.Time
	if sub.EndDate != nil {
		endTime, err := time.Parse("01-2006", *sub.EndDate)
		if err != nil {
			http.Error(w, "incorrect end_date format, expected: MM-YYYY or null", http.StatusBadRequest)
			return
		}
		endDate = &endTime
	}

	q := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id uuid.UUID
	err = h.DB.QueryRow(q, sub.ServiceName, sub.Price, sub.UserID, startDate, &endDate).Scan(&id)
	if err != nil {
		log.Println("DB insert error:", err)
		http.Error(w, "unable to create subscription", http.StatusInternalServerError)
		return
	}

	sub.ID = id
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(sub)
	if err != nil {
		log.Println("Create JSON response error: ", err)
		http.Error(w, "unable to create response", http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("incorrect UUID", err)
		http.Error(w, "incorrect UUID in request", http.StatusBadRequest)
		return
	}

	sub := Subscription{}

	q := `SELECT id, service_name, price, user_id, start_date, end_date
	          FROM subscriptions WHERE id = $1`

	row := h.DB.QueryRow(q, id)
	var startDate time.Time
	var endDate sql.NullTime

	err = row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &startDate, &endDate)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("subscription with uuid %s not found in DB", id)
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	sub.StartDate = startDate.Format("01-2006")

	if endDate.Valid {
		t := endDate.Time
		str := t.Format("01-2006")
		sub.EndDate = &str
	}

	responseBytes, err := json.Marshal(sub)
	if err != nil {
		log.Println("Unable converting response to JSON")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)
}
