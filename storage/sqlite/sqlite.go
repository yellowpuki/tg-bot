// Package sqlite ...
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yellowpuki/tg-bot/storage"
)

// Storage a type of sql storage.
type Storage struct {
	db *sql.DB
}

// New create a new SQL Storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

// Save save a page in storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages(url, user_name, created) VALUES(?, ?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName, time.Now().String()); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// TODO: implement PickLast
// PickLast pick a last page from storage
func (s *Storage) PickLast(ctx context.Context, userNmae string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY created DESC LIMIT 1`

	var url string
	err := s.db.QueryRowContext(ctx, q, userNmae).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick last page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userNmae,
	}, nil
}

// PickRandom pick a random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name=? ORDER BY RANDOM() LIMIT 1`

	var url string
	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove remove a page from storage.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url=? AND user_name=?`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't delete page: %w", err)
	}

	return nil
}

// IsExists cheks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

// Init create a table in SQL storage if it has not been yet.
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT, created TEXT)`
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
