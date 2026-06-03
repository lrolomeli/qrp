package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("set wal: %w", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS clicks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL,
			destination_url TEXT NOT NULL,
			ip TEXT,
			user_agent TEXT,
			referer TEXT,
			created_at TEXT NOT NULL DEFAULT (datetime('now'))
		)
	`); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) RecordClick(code, destURL, ip, ua, referer string) error {
	_, err := s.db.Exec(
		"INSERT INTO clicks (code, destination_url, ip, user_agent, referer) VALUES (?, ?, ?, ?, ?)",
		code, destURL, ip, ua, referer,
	)
	return err
}

func (s *Store) GetTotalClicks() (int, error) {
	var total int
	err := s.db.QueryRow("SELECT COUNT(*) FROM clicks").Scan(&total)
	return total, err
}

func (s *Store) GetClicksByCode() ([]ClickCount, error) {
	rows, err := s.db.Query(`
		SELECT code, destination_url, COUNT(*) as clicks
		FROM clicks
		GROUP BY code
		ORDER BY clicks DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []ClickCount
	for rows.Next() {
		var c ClickCount
		if err := rows.Scan(&c.Code, &c.DestinationURL, &c.Clicks); err != nil {
			return nil, err
		}
		counts = append(counts, c)
	}
	return counts, rows.Err()
}

func (s *Store) GetRecentClicks(limit int) ([]Click, error) {
	rows, err := s.db.Query(`
		SELECT id, code, destination_url,
		       COALESCE(ip,''), COALESCE(user_agent,''),
		       COALESCE(referer,''), created_at
		FROM clicks
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []Click
	for rows.Next() {
		var c Click
		if err := rows.Scan(&c.ID, &c.Code, &c.DestinationURL,
			&c.IP, &c.UserAgent, &c.Referer, &c.CreatedAt); err != nil {
			return nil, err
		}
		clicks = append(clicks, c)
	}
	return clicks, rows.Err()
}


