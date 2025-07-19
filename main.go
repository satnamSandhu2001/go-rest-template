package main

import (
	"context"
	"flag"
	"fmt"

	"go-rest-template/internal/db"
	dbConn "go-rest-template/internal/db/conn"
	"go-rest-template/internal/routes"
	"os"

	"go-rest-template/pkg/config"
	"go-rest-template/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.InfoF("Running in `%s` mode", config.APP().GO_ENV)
	// Command line flags
	var (
		runMigrations = flag.Bool("migrate", false, "Run database migrations")
		rollback      = flag.Int("rollback", 0, "Rollback n migrations")
		showVersion   = flag.Bool("version", false, "Show current migration version")
	)
	flag.Parse()

	// Database connection
	pool := dbConn.ConnectToDB()
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		logger.Panic("Failed to ping database:", err)
	}

	// migrations
	migrationPath := "./schema/migrations"

	switch {
	case *runMigrations || config.APP().MIGRATE_ON_START:
		logger.Info("Running migrations...")
		if err := dbConn.RunMigrations(pool, migrationPath); err != nil {
			logger.PanicF("Migration failed: %v", err)
		}
		logger.Info("Migration completed.")
		if !config.APP().MIGRATE_ON_START {
			os.Exit(0)
		}
	case *rollback > 0:
		logger.InfoF("Rolling back %d migrations...\n", *rollback)
		if err := dbConn.RollbackMigrations(pool, migrationPath, *rollback); err != nil {
			logger.PanicF("Rollback failed: %v", err)
		}
		os.Exit(0)
	case *showVersion:
		v, dirty, err := dbConn.GetMigrationVersion(pool)
		if err != nil {
			logger.PanicF("Error reading migration version: %v", err)
		}
		logger.InfoF("Current migration version: %s (dirty: %v)", v, dirty)
		os.Exit(0)
	default:
		logger.Info("No migration action specified.")
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.RequestID)
	if config.APP().GO_ENV == "production" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{
				config.APP().CLIENT_URL,
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	} else {
		r.Use(devCors)
	}
	r.Use(chiMiddleware.Heartbeat("/api/health"))

	// Injectors
	q := db.New(pool)

	// register routes
	r.Route("/api", func(api chi.Router) {
		routes.RegisterUserRoutes(api, q)
	})

	logger.InfoF("Server listening on :%d", config.APP().PORT)
	logger.Panic(http.ListenAndServe(fmt.Sprintf(":%d", (config.APP().PORT)), r))
}

// CORS middleware for development environment
func devCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
