package db

import (
	"context"
	"errors"
	"fmt"
	"go-rest-template/pkg/logger"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	Version string
	UpSQL   string
	DownSQL string
}

func ensureMigrationsTable(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT now()
		);
	`)
	return err
}

func loadMigrations(dir string) ([]Migration, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	mMap := map[string]*Migration{}
	for _, f := range files {
		name := f.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		parts := strings.Split(name, ".")
		if len(parts) < 3 {
			continue
		}
		version := strings.Split(parts[0], "_")[0]
		mig, ok := mMap[version]
		if !ok {
			mig = &Migration{Version: version}
			mMap[version] = mig
		}
		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(name, ".up.sql") {
			mig.UpSQL = string(content)
		} else if strings.HasSuffix(name, ".down.sql") {
			mig.DownSQL = string(content)
		}
	}

	var migrations []Migration
	for _, m := range mMap {
		migrations = append(migrations, *m)
	}
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	return migrations, nil
}

func RunMigrations(pool *pgxpool.Pool, dir string) error {
	ctx := context.Background()
	if err := ensureMigrationsTable(pool); err != nil {
		return err
	}
	migs, err := loadMigrations(dir)
	if err != nil {
		return err
	}

	applied := map[string]bool{}
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var v string
		_ = rows.Scan(&v)
		applied[v] = true
	}

	for _, m := range migs {
		if applied[m.Version] {
			continue
		}
		logger.Info("Applying migration: %s", m.Version)
		if _, err := pool.Exec(ctx, m.UpSQL); err != nil {
			return fmt.Errorf("error applying %s: %w", m.Version, err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, m.Version); err != nil {
			return err
		}
	}
	return nil
}

func RollbackMigrations(pool *pgxpool.Pool, dir string, n int) error {
	ctx := context.Background()
	if err := ensureMigrationsTable(pool); err != nil {
		return err
	}
	migs, err := loadMigrations(dir)
	if err != nil {
		return err
	}

	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT $1`, n)
	if err != nil {
		return err
	}
	defer rows.Close()

	var toRollback []string
	for rows.Next() {
		var v string
		_ = rows.Scan(&v)
		toRollback = append(toRollback, v)
	}

	if len(toRollback) == 0 {
		return errors.New("no migrations to rollback")
	}

	for _, rev := range toRollback {
		for _, m := range migs {
			if m.Version == rev {
				logger.Info("Rolling back migration: %s", m.Version)
				if _, err := pool.Exec(ctx, m.DownSQL); err != nil {
					return fmt.Errorf("rollback error %s: %w", m.Version, err)
				}
				if _, err := pool.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, m.Version); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

func GetMigrationVersion(pool *pgxpool.Pool) (string, bool, error) {
	ctx := context.Background()
	var version string
	err := pool.QueryRow(ctx, `SELECT version FROM schema_migrations ORDER BY applied_at DESC LIMIT 1`).Scan(&version)
	if err != nil {
		return "0", false, nil
	}
	return version, false, nil
}
