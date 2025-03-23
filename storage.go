package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*TransferRequest) error
	GetAccounts() ([]*Account, error)
	GetAccountByNumber(int) (*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists accounts(
    	id serial primary key,
    	first_name varchar(100) not null,
    	last_name varchar(100) not null,
    	number serial,
    	balance float,
    	created_at timestamp 	
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateAccount(a *Account) error {
	statement, err := s.db.Prepare(`
		insert into accounts(first_name, last_name, number, balance, created_at)
		values ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(a.FirstName, a.LastName, a.Number, a.Balance, a.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("delete from accounts where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) UpdateAccount(request *TransferRequest) error {
	statement, err := s.db.Prepare(`update accounts set balance = $1 where id = $2`)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(request.Amount, request.ToAccount)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account

	for rows.Next() {
		var account Account
		if err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("select * from accounts where number = $1", number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var account Account

	for rows.Next() {
		if err = rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt); err != nil {
			return nil, err
		}
	}

	return &account, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from accounts where id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var account Account

	for rows.Next() {
		if err = rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt); err != nil {
			return nil, err
		}
	}

	return &account, nil
}
