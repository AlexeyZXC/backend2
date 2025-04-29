package main

import (
	"fmt"
	"log"
	"net/http"

	"serviceReqHandler/red"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	measurable = red.MeasurableHandler
	router     = mux.NewRouter()
	web        = http.Server{
		Handler: router,
		Addr:    ":8082",
	}
)

func init() {
	router.
		HandleFunc("/identity", measurable(GetIdentityHandler)).
		Methods(http.MethodGet)
}

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9092", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}
func GetIdentityHandler(w http.ResponseWriter, r *http.Request) {
	defer log.Printf("GetIdentityHandler end\n")
	log.Printf("GetIdentityHandler start\n")

	if r.FormValue("token") == "admin_secret_token" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("GetIdentityHandler: ok \n")))
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf("GetIdentityHandler: StatusUnauthorized \n")))
}
