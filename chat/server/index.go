package server

import (
	"fmt"
	"net/http"

	"github.com/JakubC-projects/chat/chat"
	"github.com/rs/cors"
)

type Server struct {
	port      string
	store     chat.ChatStore
	publisher chat.Publisher
	publicDir string
}

func New(port string, store chat.ChatStore, publisher chat.Publisher, publicDir string) *Server {

	return &Server{
		port:      port,
		store:     store,
		publisher: publisher,
		publicDir: publicDir,
	}
}

func (s *Server) Run() error {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/messages", s.handleSendMessage)
	mux.HandleFunc("GET /api/events/messages", s.handleMessageEvents)
	fmt.Println("public dir", s.publicDir)
	mux.Handle("GET /", http.FileServer(http.Dir(s.publicDir)))

	handler := cors.AllowAll().Handler(mux)

	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), handler)
}
