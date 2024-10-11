package main

import (
	"academy-adventure-game/model"
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"net/http"
)

var game *model.Game

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, this is the Academy adventure game!")
}

func startGame(writer http.ResponseWriter, request *http.Request) {
	game := &model.Game{}
	game.SetupGame()

	for {
		var playerInput model.PlayerInput
		if err := json.NewDecoder(request.Body).Decode(&playerInput); err != nil {
			http.Error(writer, "Invalid input", http.StatusBadRequest)
			return
		}

		response := game.RunGame(playerInput)

		if err := json.NewEncoder(writer).Encode(response); err != nil {
			http.Error(writer, "Error encoding response", http.StatusInternalServerError)
			return
		}

		if response.GameOver {
			break
		}
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(writer, request)
	})
}

func main() {

	router := http.NewServeMux()

	router.HandleFunc("GET /", rootHandler)
	router.HandleFunc("GET /NewGame", startGame)

	game = &model.Game{}
	game.SetupGame()

	// commands := []string{"start", "look"}

	// for _, command := range commands {

	// 	gameResponse := game.RunGame(model.PlayerInput{Command: command, Args: []string{}})

	// 	fmt.Printf("response: %s", gameResponse.Message)
	// }

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
