package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq" // for database/sql

	"github.com/ilyakharev/url-short/internal/storage"
)

type Storage struct {
	db *sql.DB
}

const (
	templateTable = `
CREATE TABLE IF NOT EXISTS urls (
	short_url	VARCHAR(10) PRIMARY KEY,
	full_url    VARCHAR(1024)
);

CREATE INDEX IF NOT EXISTS idx ON urls USING hash(
	full_url
);
`
	templateGetFullURL  = `SELECT full_url FROM urls WHERE short_url = $1`
	templateInsertShort = `INSERT INTO urls(short_url, full_url) VALUES ($1, $2)`
	templateCheckExists = `SELECT short_url FROM urls WHERE full_url = $1`
)

var _ storage.Storager = &Storage{}

func New(url string) (*Storage, error) {
	var err error
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(templateTable)
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (st *Storage) GetFullURL(_ context.Context,
	token string,
) (fullURL string, found bool, err error) {
	rows, err := st.db.Query(templateGetFullURL, token)
	if err != nil {
		return "", false, err
	}
	if rows.Err() != nil {
		return "", false, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			return
		}
	}()
	if !rows.Next() {
		return "", false, nil
	}
	err = rows.Scan(&fullURL)
	if err != nil {
		return "", false, err
	}
	return fullURL, true, nil
}

func (st *Storage) CreateShortURL(_ context.Context, fullURL string,
	token string,
) (err error) {
	_, err = st.db.Exec(templateInsertShort,
		token, fullURL)
	return err
}

func (st *Storage) AlreadyExists(_ context.Context,
	fullURL string,
) (token string, found bool, err error) {
	rows, err := st.db.Query(templateCheckExists, fullURL)
	if err != nil {
		return "", false, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			return
		}
	}()

	if rows.Err() != nil {
		return "", false, err
	}
	if !rows.Next() {
		return "", false, nil
	}
	err = rows.Scan(&token)
	if err != nil {
		return "", false, err
	}
	return token, true, nil
}

func (st *Storage) Close() error {
	return st.db.Close()
}
