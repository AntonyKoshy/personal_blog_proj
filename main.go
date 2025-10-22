package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// Define constants for your DB connection parameters
const (
	dbUser = "antonykoshy"
	dbHost = "localhost"
	dbPort = 5432
	dbName = "blogdb"
)

// Struct representing articles table
type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Category  string    `json:"catergory"`
	Tags      []string  `json:"tags"`
}

// jwt secret key: to be set in env variable
var jwtKey = []byte("your_secret_key_here")

var conn *pgx.Conn

func main() {

	dsn := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)

	//Connect to db
	var err error
	conn, err = pgx.Connect(context.Background(), dsn)

	if err != nil {
		log.Fatalf("Unable to connect to database :%v\n", err)

	}

	defer conn.Close(context.Background())

	//test query for testing db connection
	var greeting string
	err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&greeting)
	if err != nil {
		log.Fatalf("Failed to execute test query: %v\n", err)
	}
	fmt.Println("Database connection verified")

	//Registering the endpoint
	http.HandleFunc("/articles", getArticlesHandler)

	//Starting the server
	fmt.Println("Starting server on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}

func generateJWT(username string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "kundan blogs",
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Usually token comes as "Bearer <token>", so parse it
		parts := strings.Split(tokenStr, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenStr = parts[1]

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Query the database for articles
		rows, err := conn.Query(context.Background(), "SELECT id, content, title , created_at, updated_at, category, tags FROM articles")
		if err != nil {
			http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var articles []Article

		for rows.Next() {
			var a Article
			if err := rows.Scan(&a.ID, &a.Content, &a.Title, &a.CreatedAt, &a.UpdatedAt, &a.Category, &a.Tags); err != nil {
				http.Error(w, "Error scanning articles", http.StatusInternalServerError)
				return
			}
			articles = append(articles, a)
		}

		// Return the list as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(articles)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}
