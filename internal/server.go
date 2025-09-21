package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func renderJSON[T ResponseError | ResponseSuccess](w http.ResponseWriter, response T) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

type HttpServer struct {
	mux *http.ServeMux
}

// sets the route handlers

func (server *HttpServer) initHandlers() {
	server.mux.HandleFunc("/get/{file}/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("file")
		renderJSON(w, ResponseSuccess{"hello world: " + fileName})
	})

	server.mux.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		renderJSON(w, NewResponse400())
	})

	server.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderJSON(w, NewResponse404())
	})
}

// creates the ServeMux, inits the route handlers, and starts the server

func (server *HttpServer) StartServer() {
	server.mux = http.NewServeMux()

	server.initHandlers()

	log.Print("Server started...")

	log.Fatal(http.ListenAndServe(":8000", server.mux))

}
