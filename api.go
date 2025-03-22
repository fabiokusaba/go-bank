package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Estrutura da nossa API Server
type APIServer struct {
	listenAddr string
	store      Storage
}

// Criando uma nova instância de APIServer
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Inicializando a nossa API Server
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))

	log.Println("API Server running on port:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

// Handlers responsáveis por executar determinadas ações
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// Opção com switch case
	switch r.Method {
	case http.MethodGet:
		return s.handleGetAccount(w, r)
	case http.MethodPost:
		return s.handleCreateAccount(w, r)
	case http.MethodDelete:
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("Method not allowed %s", r.Method)
	}

	// Opção com if statements
	// if r.Method == http.MethodGet {
	// 	return s.handleGetAccount(w, r)
	// }
	// if r.Method == http.MethodPost {
	// 	return s.handleCreateAccount(w, r)
	// }
	// if r.Method == http.MethodDelete {
	// 	return s.handleDeleteAccount(w, r)
	// }
	// return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
