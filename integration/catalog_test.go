package integration_test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/theguuholi/rent-cars-go/internal/db/store"
	"github.com/theguuholi/rent-cars-go/migrations"
	"slices"
	"testing"
	"time"
)

func TestCatalogOnPostgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	container, err := postgres.Run(ctx, "postgres:17", postgres.WithDatabase("catalog"), postgres.WithUsername("test"), postgres.WithPassword("test"), postgres.BasicWaitStrategies())
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatal(err)
	}
	url, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()
	for range 2 {
		if err := migrations.Up(ctx, db); err != nil {
			t.Fatal(err)
		}
	}
	q := store.New(pool)
	rows, err := q.ListCategories(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, r := range rows {
		names = append(names, r.Name)
	}
	if !slices.Equal(names, []string{"HATCHBACK", "PICKUP", "SEDAN", "SUV"}) {
		t.Fatalf("unexpected seed catalog: %v", names)
	}
	if _, err = pool.Exec(ctx, "INSERT INTO categories (name) VALUES ('VAN')"); err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, "INSERT INTO categories (name) VALUES ('VAN')"); err == nil {
		t.Fatal("database accepted duplicate name")
	}
	rows, err = q.ListCategories(ctx)
	if err != nil {
		t.Fatal(err)
	}
	last := rows[len(rows)-1]
	if last.Name != "VAN" || last.Description != nil {
		t.Fatalf("unexpected nullable category: %+v", last)
	}
}
