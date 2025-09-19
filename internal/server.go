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
	Mux *http.ServeMux
}

// sets the route handlers

func (server *HttpServer) initHandlers() {
	server.Mux.HandleFunc("/get/{file}/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("file")
		renderJSON(w, ResponseSuccess{"hello world: " + fileName})
	})

	server.Mux.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		renderJSON(w, NewResponse400())
	})

	server.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderJSON(w, NewResponse404())
	})
}

// creates the ServeMux, inits the route handlers, and starts the server

func (server *HttpServer) StartServer() {
	server.Mux = http.NewServeMux()

	server.initHandlers()

	log.Print("Server started...")
	log.Fatal(http.ListenAndServe(":8000", server.Mux))
}
