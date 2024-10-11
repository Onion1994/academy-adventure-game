package main

import (
	"academy-adventure-game/model"
	"fmt"
	"github.com/rs/cors"
	"net/http"
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, this is the Academy adventure game!")
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(writer, request)
	})
}

func main() {

	game := &model.Game{}
	game.SetupGame()
	// game.RunGame()

	router := http.NewServeMux()

	router.HandleFunc("GET /", rootHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(router)

	fmt.Println("Server listening on port 8080...")
    err := http.ListenAndServe(":8080", handler)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }

}

