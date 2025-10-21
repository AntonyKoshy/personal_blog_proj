package main

import (
	"time"
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// Define constants for your DB connection parameters
const (
    dbUser     = "navyapillai"
    dbHost     = "localhost"
    dbPort     = 5432
    dbName     = "blogdb"
)

// Struct representing articles table
type Article struct {
    ID        int       `json:"id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

//jwt secret key: to be set in env variable
var jwtKey = []byte("your_secret_key_here")



func main() {

	dsn := fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable",dbUser, dbHost, dbPort, dbName)
	
	//Connect to db
	conn, err := pgx.Connect(context.Background(),dsn)

	if err!=nil{
		log.Fatalf("Unable to connect to database :%v\n",err)

	}

	defer conn.Close(context.Background())
	

	//test query for testing db connection
	var greeting string
	err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&greeting)
	if err != nil {
		log.Fatalf("Failed to execute test query: %v\n", err)
	}
	fmt.Println("Database connection verified")



}

func generateJWT(username string) (string, error){
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


