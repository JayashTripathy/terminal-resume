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
	router.HandleFunc("/account", utils.MakeHTTPHandler(s.handleAccount))
	router.HandleFunc("/login", utils.MakeHTTPHandler(controllers.Login))
	utils.HttpHandlerFunc("/signup", controllers.Signup, router)
	
	fmt.Println("Starting server on", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
