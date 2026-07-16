package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var httpServer *http.Server

func runHttpRest(port string, kafka *KafkaManager, handleServerStopped func()) {
	http.HandleFunc("GET /api/events/health", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		setResponse(w, http.StatusOK, HealthResponse{Status: true})
	})
	http.HandleFunc("POST /api/events/user", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		handlePostImpl(w, r, "user", validateUser, kafka, TopicUsers)
	})
	http.HandleFunc("POST /api/events/movie", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		handlePostImpl(w, r, "movie", validateMovie, kafka, TopicMovies)
	})
	http.HandleFunc("POST /api/events/payment", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		handlePostImpl(w, r, "payment", validatePayment, kafka, TopicPayments)
	})

	httpServer = &http.Server{
		Addr: ":" + port,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка HTTP сервера: %v", err)
		}
		handleServerStopped()
	}()
}

func shutdownHttpRest() {
	if httpServer == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке HTTP сервера: %v", err)
	}
}

func handlePostImpl[T any](w http.ResponseWriter, r *http.Request, objType string, val func(T) error, kafka *KafkaManager, topic string) {
	req, err := parseRequestJsonBody(r, val)
	if err != nil {
		setError(w, http.StatusBadRequest, err.Error())
		return
	}

	ev := createEvent(objType, req)

	partition, offset, err := kafka.Produce(topic, ev)
	if err != nil {
		setError(w, http.StatusInternalServerError, "Не удалось опубликовать событие в Kafka: "+err.Error())
		return
	}

	response := EventResponse{
		Status:    "success",
		Partition: partition,
		Offset:    offset,
		Event:     ev,
	}

	setResponse(w, http.StatusCreated, response)
}

func logRequest(r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
}

func parseRequestJsonBody[T any](r *http.Request, val func(T) error) (T, error) {
	defer r.Body.Close()

	var obj T

	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		return obj, fmt.Errorf("Содержимое запроса должно быть в формате JSON - %w", err)
	}
	err = val(obj)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func validateUser(u UserEvent) error {
	if u.UserId <= 0 {
		return fmt.Errorf("Поле user_id обязательное")
	}
	if u.Action == "" {
		return fmt.Errorf("Поле action обязательное")
	}
	if u.Timestamp.IsZero() {
		return fmt.Errorf("Поле timestamp обязательное")
	}
	return nil
}

func validateMovie(m MovieEvent) error {
	if m.MovieId <= 0 {
		return fmt.Errorf("Поле movie_id обязательное")
	}
	if m.Title == "" {
		return fmt.Errorf("Поле title обязательное")
	}
	if m.Action == "" {
		return fmt.Errorf("Поле action обязательное")
	}
	return nil
}

func validatePayment(p PaymentEvent) error {
	if p.PaymentId <= 0 {
		return fmt.Errorf("Поле payment_id обязательное")
	}
	if p.UserId <= 0 {
		return fmt.Errorf("Поле user_id обязательное")
	}
	if p.Amount <= 0 {
		return fmt.Errorf("Поле amount должно быть положительным")
	}
	if p.Status == "" {
		return fmt.Errorf("Поле status обязательное")
	}
	if p.Timestamp.IsZero() {
		return fmt.Errorf("Поле timestamp обязательное")
	}
	return nil
}

func createEvent(objType string, payload any) Event {
	return Event{
		Id:        uuid.New().String(),
		Type:      objType,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

func setError(w http.ResponseWriter, status int, message string) {
	setResponse(w, status, Error{Error: message})
}

func setResponse(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("При отправке результата по HTTP произошла ошибка: %v", err)
	}
}
