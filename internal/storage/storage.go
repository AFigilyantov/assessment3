package storage

import (
	"chitests/internal/models"
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type SQLiteStorage struct {
	db *sql.DB
}

func New(dbPath string) (SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return SQLiteStorage{}, err
	}
	state, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS users
		(
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS ids_username ON users (
     										username
                                            );
		`)

	if err != nil {
		return SQLiteStorage{}, err
	}
	_, err = state.Exec()
	if err != nil {
		return SQLiteStorage{}, err
	}
	return SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) RegisterUser(ctx context.Context, u models.UserAccount) (models.UserAccount, error) {

	hashedPasssword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {

		return models.UserAccount{}, err

	}

	st, err := s.db.PrepareContext(ctx, "INSERT INTO users (username, password, email) VALUES (?,?,?);")

	if err != nil {
		return models.UserAccount{}, err
	}

	if _, err := st.Exec(u.UserName, string(hashedPasssword[:]), u.Email); err != nil {
		return models.UserAccount{}, err
	}

	return models.UserAccount{}, nil
}

func (s *SQLiteStorage) LoginUser(ctx context.Context, username, password string) (models.UserAccount, error) {

	passwordFromDb := ""

	if err := s.db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&passwordFromDb); err != nil {
		if err == sql.ErrNoRows {
			return models.UserAccount{}, err
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(passwordFromDb), []byte(password))
	if err != nil {
		log.Println("Error in storage module")
		return models.UserAccount{}, err
	}

	return models.UserAccount{UserName: username, Password: password}, nil
}

func (s *SQLiteStorage) CloseDb() {
	s.db.Close()
}
