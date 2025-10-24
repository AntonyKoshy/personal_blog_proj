package main

import (
	"context"
	"log"
	"net/http"
	"personal-blog/internals/config"
	"personal-blog/internals/handlers"
	"personal-blog/internals/repository"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define constants for your DB connection parameters
const (
	dbUser = "antonykoshy"
	dbHost = "localhost"
	dbPort = 5432
	dbName = "blogdb"
)

// jwt secret key: to be set in env variable
var jwtKey = []byte("your_secret_key_here")

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Unable to fetch env variables: %v\n", err)
	}

	db, err := repository.ConnectDB(cfg)

	if err != nil {
		log.Fatalf("Unable to connect to DB: %v\n", err)

	}

	defer db.Close(context.Background())

	articleRepo := &repository.ArticleRepository{Conn: db}

	// Setup HTTP router P
	mux := http.NewServeMux()

	// Register protected routes
	// mux.Handle("/articles", http.HandlerFunc(handlers.GetArticlesHandler(articleRepo)))
	// mux.Handle("/articles", http.HandlerFunc(handlers.CreateArticleHandler(articleRepo)))
	mux.HandleFunc("GET /articles", handlers.GetArticlesHandler(articleRepo))
	mux.HandleFunc("GET /articles/{id}", handlers.GetArticlesByIDHandler(articleRepo))

	mux.HandleFunc("POST /articles", handlers.CreateArticleHandler(articleRepo))

	// Start server
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

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
