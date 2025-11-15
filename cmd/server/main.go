package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	userpb "github.com/nicolasmaurizi/go-grpc-rest-basics/proto"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Implementación del servicio gRPC
type userServer struct {
	userpb.UnimplementedUserServiceServer
	db *sql.DB
}

func (s *userServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	log.Printf("CreateUser recibido: name=%s, email=%s", req.GetName(), req.GetEmail())
	log.Printf("RAW request (Go struct): %#v", req)

	name := req.GetName()
	email := req.GetEmail()

	var id int64
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		name, email,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		User: &userpb.User{
			Id:    id,
			Name:  name,
			Email: email,
		},
	}, nil
}

// Handler REST para listar usuarios
func listUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT id, name, email FROM users ORDER BY id`)
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			log.Println("query error:", err)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				http.Error(w, "scan error", http.StatusInternalServerError)
				log.Println("scan error:", err)
				return
			}
			users = append(users, u)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "json error", http.StatusInternalServerError)
			return
		}
	}
}

func logRawInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Log visible del método
	log.Printf("→ gRPC call: %s", info.FullMethod)

	// request = Protobuf => raw bytes
	if msg, ok := req.(proto.Message); ok {
		raw, _ := proto.Marshal(msg)
		log.Printf("↳ RAW BYTES: %v", raw)
		log.Printf("↳ HEX: %X", raw)
		log.Printf("↳ STRUCT: %#v", msg)
	}
	return handler(ctx, req)
}

func main() {
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

	// ---- Servidor gRPC ----
	grpcLis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen grpc: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logRawInterceptor),
	)
	userpb.RegisterUserServiceServer(grpcServer, &userServer{db: db})

	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
		}
	}()

	// ---- Servidor REST ----
	mux := http.NewServeMux()
	mux.Handle("/users", listUsersHandler(db))

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("HTTP server listening on :8080")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("http server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
