package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	// serveMux is a serve multiplexer that matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.
	serveMux := http.NewServeMux()

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	serveMux.Handle("/app/",
		http.StripPrefix("/app/",
			apiCfg.middlewareMetricsInc(
				http.FileServer(http.Dir("./static")),
			),
		),
	)

	// Custom handler for the "/healthz" endpoint that responds with a 200 OK status and a plain text message.
	serveMux.HandleFunc("/healthz", handlerHealth)

	serveMux.HandleFunc("/metrics", apiCfg.handlerHits)
	serveMux.HandleFunc("/reset", apiCfg.handlerReset)

	// log.Fatal is used to log any error that occurs when starting the server. Triggered when the server fails to start or encounters an error while running.
	log.Fatal(server.ListenAndServe())
}

// custom handler function for the "/healthz" endpoint.
func handlerHealth(w http.ResponseWriter, req *http.Request) {
	// .Set vs .Add : Set replaces the value of the header with the provided value, while Add appends the value to the existing values for that header. Set is better since we want to ensure that the Content-Type is set correctly without duplicates.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// http.StatusOK is used to match the response body match with status code.
	// Ex: if w.WriteHeader(http.StatusServiceUnavailable) is used, will result in a 503 Service Unavailable status code
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	prev := cfg.fileserverHits.Load() // get current value
	cfg.fileserverHits.Store(0)       // reset to 0

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits before reset: %d\n", prev)
}
