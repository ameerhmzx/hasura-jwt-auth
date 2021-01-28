package main

import (
	"actions/auth"
	"actions/utils"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func runServer(port string, debug bool) {
	router := mux.NewRouter()

	// Handlers
	auth.Handler(router)

	if debug {
		router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("Actions Working ...\n"))
		})
		fmt.Println("Listening on " + port)
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func main() {
	runServer(
		utils.LoadEnv("PORT", "8081"),
		utils.LoadEnvBool("DEBUG", false),
		)
}
