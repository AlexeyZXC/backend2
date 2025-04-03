package main

import (
	"L1/internal/api/router"
	"L1/internal/app"
	"L1/internal/db"
	"L1/internal/domain"
)

func main() {

	db := db.NewDB()
	a := app.NewApp(db)
	service := domain.NewService(a)
	router := router.NewRouter(service)
	router.Start()

}
