package main

import "math/rand"

// Definindo a estrutura de uma conta
type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

// Criando uma nova instância de Account
func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(1000000)),
		// Não precisamos passar o Balance porque o Go automaticamente vai inicializar com 0, valor default para int64
	}
}
