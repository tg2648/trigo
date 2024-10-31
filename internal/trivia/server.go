package trivia

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	staticDir string
	hub       *Hub
	mux       *http.ServeMux
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewServer(staticDir string, hub *Hub) *Server {
	mux := http.NewServeMux()
	server := &Server{staticDir: staticDir, hub: hub, mux: mux}
	server.configureRoutes()

	return server
}

func (s *Server) configureRoutes() {
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("GET /room/{room}", s.apiHandleRoomGet)

	rootMux := http.NewServeMux()
	rootMux.Handle("/", http.FileServer(http.Dir(s.staticDir)))
	rootMux.HandleFunc("GET /ws", s.handleWebsocket)

	s.mux.Handle("/", rootMux)
	s.mux.Handle("/api/", http.StripPrefix("/api", apiMux))
}

func (s *Server) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: s.hub, send: make(chan []byte, 256), conn: conn}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
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
