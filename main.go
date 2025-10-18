package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"fmt"
)

// Define constants for your DB connection parameters
const (
    dbUser     = "navyapillai"
    dbHost     = "localhost"
    dbPort     = 5432
    dbName     = "postgres"
)

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