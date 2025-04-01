package main

import (
	"L1/internal/app"
	"L1/internal/db"
)

func main() {

	db := db.NewDB()

	a := app.NewApp(db)

}
