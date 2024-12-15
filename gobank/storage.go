package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (int, error)
	GetAccounts() ([]*Account, error)
	GetAccountByID(id int) (*Account, error)
	UpdateAccount(account *Account) error
	DeleteAccount(id int) int
	TransferFunds(fromId int, toId int, amount int64) error
}

type PostgresDB struct {
	db *sql.DB
}

func (store *PostgresDB) Init() error {
	// create tables if not exist
	query := `
  CREATE TABLE IF NOT EXISTS accounts( 
    id serial PRIMARY KEY,
    first_name VARCHAR(255), 
    last_name VARCHAR(255), 
    number int,
    balance int, 
    created_at timestamp
    );`

	_, err := store.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostgresDB) GetAccountByID(id int) (*Account, error) {
	account := &Account{ID: -1}
	row := store.db.QueryRow("select * from accounts where id = $1", id)

	err := row.Scan(&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (store *PostgresDB) UpdateAccount(account *Account) error {
	query := `
    update accounts
    set first_name = $1, last_name = $2
    where id = $3
    returning *
  `
	row := store.db.QueryRow(query, account.FirstName, account.LastName, account.ID)

	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		log.Println("Unable to update account", err)
		return err
	}

	return nil
}
func (store *PostgresDB) CreateAccount(account *Account) (int, error) {
	query := `
  INSERT INTO accounts(
    first_name, 
    last_name,
    number,
    balance,
    created_at)
    VALUES($1,$2,$3,$4,$5)
    RETURNING id
  `

	var userID int
	err := store.db.QueryRow(query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (store *PostgresDB) DeleteAccount(id int) int {
	deletedId := -1
	row := store.db.QueryRow("delete from accounts where id = $1 returning id", id)

	row.Scan(&deletedId)

	return deletedId
}

func (store *PostgresDB) TransferFunds(fromId int, toId int, amount int64) error {
	// Create query

	// Run query

	return nil
}

func (store *PostgresDB) GetAccounts() ([]*Account, error) {
	// Create query
	rows, err := store.db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}

		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)

	}

	return accounts, nil
}

func NewPostgresDB(conStr string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}
