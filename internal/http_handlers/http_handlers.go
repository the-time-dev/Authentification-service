package http_handlers

import (
	"auth-service/internal/auth"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	router     *mux.Router
	authorizer *auth.Authorizer
}

func NewServer(authorizer *auth.Authorizer) (*Server, error) {
	r := mux.NewRouter()
	server := &Server{router: r, authorizer: authorizer}
	r.HandleFunc("/access/{userID}", server.acessHandler).Methods("POST")
	r.HandleFunc("/refresh", server.refreshHandler).Methods("POST")

	return server, nil
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func (server *Server) acessHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]

	accessToken, refreshToken, err := server.authorizer.GenerateTokenPair(userID, r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating token pair"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"accessToken": "%s", "refreshToken": "%s"}`, accessToken, refreshToken)))
}

func (server *Server) refreshHandler(w http.ResponseWriter, r *http.Request) {

	accessToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	refreshToken := r.Header.Get("Refresh")

	accessToken, refreshToken, err := server.authorizer.RefreshAccessToken(accessToken, refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error refreshing access token\n" + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"accessToken": "%s", "refreshToken": "%s"}`, accessToken, refreshToken)))
}
