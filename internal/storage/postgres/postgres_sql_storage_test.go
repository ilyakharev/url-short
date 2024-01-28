package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqlStorage_AlreadyExists(t *testing.T) {
	tests := []*struct {
		name         string
		fullURL      string
		queryError   bool
		alreadyExist bool
		token        string
	}{
		{
			name:         "query error",
			fullURL:      "http://ya.ru",
			queryError:   true,
			alreadyExist: false,
			token:        "",
		},
		{
			name:         "fullURL already exists",
			fullURL:      "http://ya.ru",
			queryError:   false,
			alreadyExist: true,
			token:        "1234567890",
		},
		{
			name:         "not found",
			fullURL:      "http://ya.ru",
			queryError:   false,
			alreadyExist: false,
			token:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			ctx := context.Background()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			st := &Storage{
				db: db,
			}
			defer func() {
				err = st.Close()
				if err != nil {
					return
				}
			}()

			switch {
			case tt.queryError:
				mock.ExpectQuery("SELECT short_url").WithArgs(
					tt.fullURL).WillReturnError(errors.New("any"))
			case tt.alreadyExist:
				rows := sqlmock.NewRows([]string{"short_url"}).AddRow(tt.token)
				mock.ExpectQuery("SELECT short_url").WithArgs(
					tt.fullURL).WillReturnRows(rows)
			case !tt.alreadyExist:
				rows := sqlmock.NewRows([]string{"short_url"})
				mock.ExpectQuery("SELECT short_url").WithArgs(
					tt.fullURL).WillReturnRows(rows)
			}

			token, found, err := st.AlreadyExists(ctx, tt.fullURL)
			switch {
			case tt.queryError:
				assert.False(t, found)
				assert.Empty(t, token)
				require.Error(t, err)
			case tt.alreadyExist:
				assert.True(t, found)
				assert.Equal(t, tt.token, token)
				require.NoError(t, err)
			default:
				assert.False(t, found)
				assert.Empty(t, token)
				require.NoError(t, err)
			}
		})
	}
}

func TestSqlStorage_CreateShortURL(t *testing.T) {
	tests := []*struct {
		name       string
		fullURL    string
		queryError bool
		token      string
	}{
		{
			name:       "query error",
			fullURL:    "http://ya.ru",
			queryError: true,
			token:      "1234567890",
		},
		{
			name:       "success",
			fullURL:    "http://ya.ru",
			queryError: false,
			token:      "1234567890",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			ctx := context.Background()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			st := &Storage{
				db: db,
			}
			defer func() {
				err = st.Close()
				if err != nil {
					return
				}
			}()

			if tt.queryError {
				mock.ExpectExec("INSERT").
					WithArgs(tt.token, tt.fullURL).
					WillReturnError(errors.New("some"))
			} else {
				mock.ExpectExec("INSERT").
					WithArgs(tt.token, tt.fullURL).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = st.CreateShortURL(ctx, tt.fullURL, tt.token)
			if tt.queryError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSqlStorage_GetFullURL(t *testing.T) {
	tests := []*struct {
		name       string
		fullURL    string
		queryError bool
		found      bool
		token      string
	}{
		{
			name:       "query error",
			fullURL:    "http://ya.ru",
			queryError: true,
			found:      false,
			token:      "",
		},
		{
			name:       "found",
			fullURL:    "http://ya.ru",
			queryError: false,
			found:      true,
			token:      "1234567890",
		},
		{
			name:       "not found",
			fullURL:    "http://ya.ru",
			queryError: false,
			found:      false,
			token:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			ctx := context.Background()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			st := &Storage{
				db: db,
			}
			defer func() {
				err = st.Close()
				if err != nil {
					return
				}
			}()

			switch {
			case tt.queryError:
				mock.ExpectQuery("SELECT full_url").WithArgs(
					tt.token).WillReturnError(errors.New("any"))
			case tt.found:
				rows := sqlmock.NewRows([]string{"full_url"}).AddRow(tt.fullURL)
				mock.ExpectQuery("SELECT full_url").WithArgs(
					tt.token).WillReturnRows(rows)
			case !tt.found:
				rows := sqlmock.NewRows([]string{"full_url"})
				mock.ExpectQuery("SELECT full_url").WithArgs(
					tt.token).WillReturnRows(rows)
			}

			fullURL, found, err := st.GetFullURL(ctx, tt.token)
			switch {
			case tt.queryError:
				assert.False(t, found)
				assert.Empty(t, fullURL)
				require.Error(t, err)
			case tt.found:
				assert.True(t, found)
				assert.Equal(t, tt.fullURL, fullURL)
				require.NoError(t, err)
			case !tt.found:
				assert.False(t, found)
				assert.Empty(t, fullURL)
				require.NoError(t, err)
			}
		})
	}
}
