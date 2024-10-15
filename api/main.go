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

func (s *APIServer) Start() {
	router := mux.NewRouter()

	utils.HttpHandlerFunc("/signup", controllers.Signup, router)
	utils.HttpHandlerFunc("/login", controllers.Login, router)
	
	fmt.Println("Starting server on", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}


