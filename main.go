package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/i101dev/multimodal-db/routes"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("Invalid port - not found in environment")
	}

	// -----------------------------------------------------------------------
	// Server Setup
	//
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
	}

	// -----------------------------------------------------------------------
	// Routing Setup
	//
	routes.RegisterTestRoutes()
	routes.RegisterUserRoutes()
	routes.RegisterAlertRoutes()
	routes.RegisterTxnRoutes()

	// -----------------------------------------------------------------------
	// Server Launch
	//
	fmt.Println("Server is live on port:", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
