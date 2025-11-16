package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nicolasmaurizi/go-grpc-rest-basics/config"
)

func NewPostgres(cfg config.DBConfig) *sql.DB {
	connStr := cfg.ConnString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("error ping db: %v", err)
	}

	return db
}

/*
	// DATABASE_URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			getenv("DB_HOST", "localhost"),
			getenv("DB_PORT", "5432"),
			getenv("DB_USER", "postgres"),
			getenv("DB_PASSWORD", "admin"),
			getenv("DB_NAME", "bloomgrpc"),
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("error ping db: %v", err)
	}
*/
