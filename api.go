package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Função para envio de respostas em formato JSON
func WriteJSON(w http.ResponseWriter, status int, value any) error {
	// Configurando o cabeçalho da resposta
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	// Enviando a resposta
	return json.NewEncoder(w).Encode(value)
}

// Tipo das nossas handle functions
type apiFunc func(http.ResponseWriter, *http.Request) error

// Definindo a estrutura da nossa API Error
type APIError struct {
	Error string
}

// Convertendo a nossa handle function para um http.HandlerFunc para estar de acordo com o nosso router
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Tratar o erro
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

// Estrutura da nossa API Server
type APIServer struct {
	listenAddr string
}

// Criando uma nova instância de APIServer
func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

// Inicializando a nossa API Server
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

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
