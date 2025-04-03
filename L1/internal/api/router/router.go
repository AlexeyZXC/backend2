package router

import (
	"L1/internal/domain"
	"net/http"
	"strconv"
)

type Router struct {
	http.ServeMux
	service *domain.Service
}

func NewRouter(s *domain.Service) *Router {
	mux := &Router{service: s}
	mux.ServeMux.HandleFunc("/newuser", mux.HandleNewUser)
	mux.ServeMux.HandleFunc("/newgroup", mux.HandleNewGroup)
	mux.ServeMux.HandleFunc("/", mux.HandleIndex)
	return mux
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.ServeMux.ServeHTTP(w, req)
}

func (r *Router) HandleIndex(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hi"))
}

func (r *Router) HandleNewUser(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	user, err := r.service.NewUser(req.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(user.Name))
}

func (r *Router) HandleNewGroup(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	gtypestr := req.URL.Query().Get("gtype")
	gtype, err := strconv.Atoi(gtypestr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, err := r.service.NewGroup(req.Context(), name, domain.GroupType(gtype))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(group.Name))
}

func (r *Router) Start() {
	http.ListenAndServe(":8080", r)
}
