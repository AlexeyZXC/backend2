package api

import (
	"fmt"
	"net/http"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

// type ServerInterface interface {
// 	// Command
// 	// (POST /{cmd}/{src}/{id})
// 	Command(w http.ResponseWriter, r *http.Request, cmd string, src string, id string)
// }

func (Server) Command(w http.ResponseWriter, r *http.Request, cmd string, src string, id string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("pong %s %s %s", cmd, src, id)))
}
