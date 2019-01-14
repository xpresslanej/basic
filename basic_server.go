package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type helloWorldResponse struct {
	Message string `json:"message"`
}

/*type helloWorldRequest struct {
	Name string `json:"name"`
}*/

//GZipHandler - handler for gzip compression
type GZipHandler struct {
	next http.Handler
}

//GZipResponseWriter - writer for gzip compression
type GZipResponseWriter struct {
	gw *gzip.Writer
	http.ResponseWriter
}

func main() {
	port := 8080

	cathandler := http.FileServer(http.Dir("./images"))
	http.Handle("/cat/", http.StripPrefix("/cat/", cathandler))
	//http.HandleFunc("/helloworld", helloWorldHandler)

	http.Handle("/helloworld", NewGZipHandler(http.HandlerFunc(helloWorldHandler)))

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	/*var request helloWorldRequest
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}*/

	//response := helloWorldResponse{Message: "Hello " + request.Name}
	response := helloWorldResponse{Message: "Hello Fool"}
	encoder := json.NewEncoder(w)
	encoder.Encode(&response)
}

//NewGZipHandler - used to create new GZipHandler
func NewGZipHandler(next http.Handler) http.Handler {
	return &GZipHandler{next}
}

func (h *GZipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encodings := r.Header.Get("Accept-Encoding")

	if strings.Contains(encodings, "gzip") {
		h.serveGZip(w, r)
	}
}

func (h *GZipHandler) serveGZip(w http.ResponseWriter, r *http.Request) {
	gzw := gzip.NewWriter(w)
	defer gzw.Close()

	w.Header().Set("Content-Encoding", "gzip")
	h.next.ServeHTTP(GZipResponseWriter{gzw, w}, r)
}

func (h *GZipHandler) serveNoCompression(w http.ResponseWriter, r *http.Request) {
	h.next.ServeHTTP(w, r)
}

func (w GZipResponseWriter) Write(b []byte) (int, error) {
	if _, ok := w.Header()["Content-Type"]; !ok {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	return w.gw.Write(b)
}

//Flush - used to flush
func (w GZipResponseWriter) Flush() {
	w.gw.Flush()
	if fw, ok := w.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}
