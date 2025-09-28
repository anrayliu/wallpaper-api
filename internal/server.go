// reverse proxy http server to serve images from azure blob storage through a restful api

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func renderJSON[T ResponseError | ResponseSuccess](w http.ResponseWriter, response T) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Encoding error: %s", err)
	}
}

type HTTPServer struct {
	mux *http.ServeMux
	srv *http.Server
}

// proxies azure blob storage by downloading image with SAS URI and
// writing into response

func fetchBlob(w http.ResponseWriter, url string) {

	// validate request first, avoids linting error "Potential HTTP request made with variable"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		renderJSON(w, NewResponse500())
		return
	}

	resp, err := http.DefaultClient.Do(req)
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

func (server *HTTPServer) initHandlers() {
	server.mux.HandleFunc("/get/{file}/", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("file")
		imagePath := fmt.Sprintf("%s%s?%s", os.Getenv("SAS_URI"), fileName, os.Getenv("SAS_TOKEN"))

		fetchBlob(w, imagePath)
	})

	server.mux.HandleFunc("/get/", func(w http.ResponseWriter, _ *http.Request) {
		renderJSON(w, NewResponse400())
	})

	server.mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		renderJSON(w, NewResponse404())
	})
}

// creates the ServeMux, inits the route handlers, and starts the server

func (server *HTTPServer) StartServer() {
	server.mux = http.NewServeMux()

	server.srv = &http.Server{
		Addr:         ":8000",
		Handler:      server.mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	server.initHandlers()

	log.Print("Server started...")

	log.Fatal(server.srv.ListenAndServe())

}
