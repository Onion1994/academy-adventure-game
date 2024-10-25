package main

import (
	"academy-adventure-game/global"
	"academy-adventure-game/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
	"net/http"
)

const cookieName = "academy-adventure-session"

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func getGame(session *sessions.Session) *model.Game {
	if game, ok := session.Values["game"].(*model.Game); ok {
		return game
	}
	return nil
}

func setGame(session *sessions.Session, game *model.Game) {
	session.Values["game"] = game
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello, this is the Academy adventure game!")
}

func startGame(writer http.ResponseWriter, request *http.Request) {

	session, _ := store.Get(request, cookieName)
	var playerInput model.PlayerInput

	err := json.NewDecoder(request.Body).Decode(&playerInput)

	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(writer, "Bad Request Body", http.StatusBadRequest)
		return
	}

	game := getGame(session)

	fmt.Println("Session before saving:", session.Values)

	if game == nil {
		game = &model.Game{}
		game.SetupGame()
		setGame(session, game)
	}

	response := game.RunGame(playerInput)

	json.NewEncoder(writer).Encode(response)

	if response.GameOver {
		global.GameOver = true
		return
	}

	session.Save(request, writer)
	fmt.Println("Session after saving:", session.Values)
}

func getAvailableActions(writer http.ResponseWriter, request *http.Request) {

	var command model.GameCommand
	session, _ := store.Get(request, cookieName)

	err := json.NewDecoder(request.Body).Decode(&command)

	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(writer, "Bad Request Body", http.StatusBadRequest)
		return
	}

	game := getGame(session)
	if game == nil {
		http.Error(writer, "No game found", http.StatusNotFound)
		return
	}

	response := game.GetAvailableActions(command.Command)

	json.NewEncoder(writer).Encode(response)

}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(writer, request)
	})
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/GameResponse", startGame)
	router.HandleFunc("/CommandOptions", getAvailableActions)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowCredentials: true,
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
