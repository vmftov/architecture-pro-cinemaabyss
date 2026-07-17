package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type HealthResponse struct {
	Status bool `json:"status"`
}

type Error struct {
	Error string `json:"error"`
}

type ReverseProxyData struct {
	TargetUri *url.URL
	Proxy     *httputil.ReverseProxy
}

var monolithProxy *ReverseProxyData
var moviesServiceProxy *ReverseProxyData
var eventsServiceProxy *ReverseProxyData
var gradualMigration bool
var moviesMigrationPercent int32

func main() {
	monolithProxy = readEnvVarAndCreateProxy("MONOLITH_URL", "http://monolith:8080")
	if monolithProxy == nil {
		return
	}
	moviesServiceProxy = readEnvVarAndCreateProxy("MOVIES_SERVICE_URL", "http://movies-service:8081")
	if moviesServiceProxy == nil {
		return
	}
	eventsServiceProxy = readEnvVarAndCreateProxy("EVENTS_SERVICE_URL", "http://events-service:8082")
	if eventsServiceProxy == nil {
		return
	}
	port := readEnvVar("PORT", "8000")
	gradualMigration = readEnvVar("GRADUAL_MIGRATION", "true") == "true"
	moviesMigrationPercent = readEnvVarPercent("MOVIES_MIGRATION_PERCENT", 50)

	http.HandleFunc("/", handleHttpRequest)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске HTTP сервера: %v", err)
	}
}

func readEnvVar(name string, fallback string) string {
	val := os.Getenv(name)
	if val == "" {
		val = fallback
		log.Printf("Не задана переменная окружения %v, используется значение %v", name, val)
	}
	return val
}

func readEnvVarPercent(name string, fallback int32) int32 {
	str := os.Getenv(name)
	if str == "" {
		log.Printf("Не задана переменная окружения %v, используется значение %v", name, fallback)
		return fallback
	}
	val, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.Printf("Не удалось преобразовать значение переменной окружения %v '%v' в целое число", name, str)
		return fallback
	}
	if val < 0 || val > 100 {
		log.Printf("Значение переменной окружения %v '%v' должно находиться в дмапазоне [0; 100]", name, str)
		return fallback
	}
	return int32(val)
}

func readEnvVarAndCreateProxy(name string, fallback string) *ReverseProxyData {
	targetAddr := readEnvVar(name, fallback)
	targetUri, err := url.Parse(targetAddr)
	if err != nil {
		log.Printf("Значение переменной окружения %v '%v' должно быть URI", name, targetAddr)
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUri)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("При перенаправлении запроса к внутреннему сервису '%s' произошла ошибка: %v", r.URL.Path, err)
		setError(w, http.StatusBadGateway, fmt.Sprintf("При перенаправлении запроса к внутреннему сервису '%s' произошла ошибка: %v", r.URL.Path, err))
	}

	return &ReverseProxyData{
		TargetUri: targetUri,
		Proxy:     proxy,
	}
}

func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	path := strings.ToLower(r.URL.Path)
	if path == "/health" || path == "/health/" {
		setPlainText(w, http.StatusOK, "Strangler Fig Proxy is healthy")
	} else if isSubPath(path, "/api/movies") {
		if shouldRedirectToMoviesMicroservice() {
			redirectTo(w, r, moviesServiceProxy)
		} else {
			redirectTo(w, r, monolithProxy)
		}
	} else if isSubPath(path, "/api/users") || isSubPath(path, "/api/payments") || isSubPath(path, "/api/subscriptions") {
		redirectTo(w, r, monolithProxy)
	} else if isSubPath(path, "/api/events") {
		redirectTo(w, r, eventsServiceProxy)
	} else {
		setError(w, http.StatusNotFound, "Ресурс не найден")
	}
}

func logRequest(r *http.Request) {
	log.Printf("%s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
}

func isSubPath(path string, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func shouldRedirectToMoviesMicroservice() bool {
	return gradualMigration && (moviesMigrationPercent == 100 || rand.Float64()*100 < float64(moviesMigrationPercent))
}

func redirectTo(w http.ResponseWriter, r *http.Request, proxyData *ReverseProxyData) {
	if r.Header.Get("X-Forwarded-For") == "" {
		r.Header.Set("X-Forwarded-For", r.RemoteAddr)
	}
	r.Header.Set("X-Forwarded-Proto", "http")
	r.Header.Set("X-Forwarded-Host", r.Host)

	r.Host = proxyData.TargetUri.Host

	proxyData.Proxy.ServeHTTP(w, r)
}

func setPlainText(w http.ResponseWriter, status int, contents string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)

	_, err := fmt.Fprint(w, contents)
	if err != nil {
		log.Printf("При отправке результата по HTTP произошла ошибка: %v", err)
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
