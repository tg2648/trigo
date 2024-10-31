package trivia

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	staticFolderDir string
	mux             *http.ServeMux
}

func NewServer(staticFolderDir string) *Server {
	mux := http.NewServeMux()
	server := &Server{staticFolderDir: staticFolderDir, mux: mux}
	server.configureRoutes()

	return server
}

func (s *Server) configureRoutes() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /room/{room}", s.apiHandleRoomGet)

	rootMux := http.NewServeMux()
	rootMux.Handle("/", http.FileServer(http.Dir(s.staticFolderDir)))

	s.mux.Handle("/", rootMux)
	s.mux.Handle("/api/", http.StripPrefix("/api", apiMux))
}

func (s *Server) apiHandleRoomGet(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("room")
	if roomID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "You requested room: "+roomID)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Run(port string) {
	log.Println("Starting server")
	log.Fatal(http.ListenAndServe(":"+port, s))
}
