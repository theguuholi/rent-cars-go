// Package migrations applies the embedded schema and catalog data.
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var scripts embed.FS

func Up(ctx context.Context, db *sql.DB) error {
	p, err := goose.NewProvider(goose.DialectPostgres, db, scripts)
	if err != nil {
		return err
	}
	_, err = p.Up(ctx)
	return err
}
