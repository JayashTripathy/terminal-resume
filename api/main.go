package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"terminal-resume.jayash.space/controllers"
	"terminal-resume.jayash.space/utils"
)

type APIServer struct {
	listenAddr string
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

// CORS middleware to handle CORS headers
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight requests
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (s *APIServer) Start() {
	router := mux.NewRouter()
	router.Use(corsMiddleware)
	utils.HttpHandlerFunc("/auth/signup", controllers.Signup, router)
	utils.HttpHandlerFunc("/auth/login", controllers.Login, router)
	
	fmt.Println("Starting server on", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}


