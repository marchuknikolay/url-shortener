package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/marchuknikolay/url-shortener/internal/config/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("could not open a database %v, err - %w", path, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		url TEXT NOT NULL,
		alias TEXT NOT NULL UNIQUE);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
	`)
	if err != nil {
		return nil, fmt.Errorf("could not create a database %v, err - %w", path, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("could not execute a query %v, err - %w", path, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave, alias string) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("could prepare a query for saving url, err - %w", err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%w", storage.ErrUrlExists)
		}

		return 0, fmt.Errorf("could execute a query for saving url, err - %w", err)
	}

	return res.LastInsertId()
}

func (s *Storage) GetUrl(alias string) (string, error) {
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("could not prepare a query for getting url for the alias %v, err - %w", alias, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%w", storage.ErrUrlNotFound)
		}

		return "", fmt.Errorf("could not find a url for the alias %v, err - %w", alias, err)
	}

	return url, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=?")
	if err != nil {
		return fmt.Errorf("could not prepare a query for deleting url for the alias %v, err - %w", alias, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("could not delete a url for the alias %v, err - %w", alias, err)
	}

	return nil
}
