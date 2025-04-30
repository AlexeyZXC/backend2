//go:generate go tool oapi-codegen -config cfg.yaml api.oapi3.yaml

///go:generate go tool oapi-codegen -generate chi-server,spec -package api -o ./api/api.go ./api.oapi3.yaml

// $ go generate
// $ go build
// $ go test

package main

import (
	"fmt"
	"l9/api"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Hello, world.")

	server := api.NewServer()
	r := chi.NewMux()
	h := api.HandlerFromMux(server, r)

	s := http.Server{
		Handler: h,
		Addr:    ":8080",
	}

	log.Fatal(s.ListenAndServe())
}
