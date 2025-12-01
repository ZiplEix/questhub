package config

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	betterlogs "github.com/ZiplEix/better-logs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// InitLogger initializes the logger and returns the logger instance and a cleanup function.
func InitLogger() (*zap.Logger, func(context.Context) error) {
	logCfg := betterlogs.DefaultConfig()
	logCfg.ServiceName = "questhub"

	var logDB *sql.DB

	// Connect to database for logs
	dsn := os.Getenv("POSTGRES_URL")
	if dsn != "" {
		var err error
		logDB, err = sql.Open("pgx", dsn)
		if err != nil {
			log.Fatalf("failed to open db for logger: %v", err)
		}

		// Ensure logs table exists
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := betterlogs.EnsureLogsTable(ctx, logDB); err != nil {
			log.Printf("warning: failed to ensure logs table: %v", err)
		}

		logCfg.DB = logDB
		logCfg.EnablePostgres = true
	}

	logger, closeLogger, err := betterlogs.New(logCfg)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	cleanup := func(ctx context.Context) error {
		if err := closeLogger(ctx); err != nil {
			return err
		}
		if logDB != nil {
			return logDB.Close()
		}
		return nil
	}

	zap.ReplaceGlobals(logger)

	return logger, cleanup
}
