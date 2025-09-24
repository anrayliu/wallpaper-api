// reverse proxy http server to serve images from azure blob storage through a restful api

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func renderJSON[T ResponseError | ResponseSuccess](w http.ResponseWriter, response T) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

type HttpServer struct {
	mux *http.ServeMux
}

// proxies azure blob storage by downloading image with SAS URI and
// writing into response

func fetchBlob(w http.ResponseWriter, url string) {

	resp, err := http.Get(url)
	if err != nil {
		renderJSON(w, NewResponse502())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		renderJSON(w, NewResponse404())
		return
	}

	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		renderJSON(w, NewResponse500())
	}

}

// sets the route handlers

func (server *HttpServer) initHandlers() {
	server.mux.HandleFunc("/get/{file}/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("file")
		imagePath := fmt.Sprintf("%s%s?%s", os.Getenv("SAS_URI"), fileName, os.Getenv("SAS_TOKEN"))

		fetchBlob(w, imagePath)
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
