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

	routes.RegisterTestRoutes()
	routes.RegisterUserRoutes()

	srv := &http.Server{
		Handler: http.DefaultServeMux,
		Addr:    ":" + port,
	}

	fmt.Println("Server is live on port:", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
