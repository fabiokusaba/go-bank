package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
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
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))

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
	case http.MethodPut:
		return s.handleTransfer(w, r)
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
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		// Pegando o id nos parâmetros da URL
		id := mux.Vars(r)["id"]
		accountID, err := strconv.Atoi(id)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(accountID)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == http.MethodDelete {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}

	account := NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := generateJWT(account)
	if err != nil {
		return err
	}

	fmt.Println("JWT Token:", tokenString)

	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(accountID); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusNoContent, nil)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}

	if err := s.store.UpdateAccount(transferRequest); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, transferRequest)
}

// Função para envio de respostas em formato JSON
func WriteJSON(w http.ResponseWriter, status int, value any) error {
	// Configurando o cabeçalho da resposta
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	// Enviando a resposta
	return json.NewEncoder(w).Encode(value)
}

const jwtSecret = "hunter999"

func generateJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission denied"})
			return
		}
		if !token.Valid {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission denied"})
			return
		}

		params := mux.Vars(r)
		accountID, err := strconv.Atoi(params["id"])
		if err != nil {
			WriteJSON(w, http.StatusNotFound, APIError{Error: "Account not found"})
			return
		}

		account, err := s.GetAccountByID(accountID)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission denied"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if account.Number != int64(claims["accountNumber"].(float64)) {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "Permission denied"})
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
}

// Tipo das nossas handle functions
type apiFunc func(http.ResponseWriter, *http.Request) error

// Definindo a estrutura da nossa API Error
type APIError struct {
	Error string `json:"error"`
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
