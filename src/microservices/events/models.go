package main

import (
	"time"
)

type HealthResponse struct {
	Status bool `json:"status"`
}

type UserEvent struct {
	UserId    int       `json:"user_id"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Username  *string   `json:"username,omitempty"`
	Email     *string   `json:"email,omitempty"`
}

type MovieEvent struct {
	MovieId     int      `json:"movie_id"`
	Title       string   `json:"title"`
	Action      string   `json:"action"`
	UserId      *int     `json:"user_id,omitempty"`
	Rating      *float64 `json:"rating,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Description *string  `json:"description,omitempty"`
}

type PaymentEvent struct {
	PaymentId  int       `json:"payment_id"`
	UserId     int       `json:"user_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
	MethodType *string   `json:"method_type,omitempty"`
}

type Event struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload"`
}

type EventResponse struct {
	Status    string `json:"status"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Event     Event  `json:"event"`
}

type Error struct {
	Error string `json:"error"`
}

const (
	TopicUsers    = "user-events"
	TopicMovies   = "movie-events"
	TopicPayments = "payment-events"
)
