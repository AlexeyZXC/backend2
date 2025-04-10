package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const version = "1.0.0"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	master := os.Getenv("PIPOMASTER") //1 - master; 0 - slave
	masterName := os.Getenv("PIPOMASTERNAME")
	slaveName := os.Getenv("PIPOSLAVENAME")
	// version := os.Getenv("PIPOVERSION")

	lg := func(format string, v ...any) {
		format = fmt.Sprintf("master(%v) version(%s) %v", master, version, format)
		log.Printf(format, v...)
	}

	defer func() {
		if r := recover(); r != nil {
			lg("Recovered %v", r)
		}
	}()

	//export PIPOMASTER=1 && export PIPOMASTERNAME=pipoMaster && export PIPOSLAVENAME=pipoSlave && export PIPOVERSION=1.0.0 && go run main.go
	//export PIPOMASTER= && export PIPOMASTERNAME= && export PIPOSLAVENAME= && export PIPOVERSION= && go run main.go

	lg("PIPOMASTER(%s) masterName(%v) slaveName(%v) version(%v)", master, masterName, slaveName, version)

	if master == "1" {
		lg("master")
		//master
		go func() {
			for {
				_, err := http.Get(fmt.Sprintf("http://%s:8080", slaveName))
				if err != nil {
					var urlErr url.Error
					if errors.Is(err, &urlErr) && urlErr.Timeout() {
						lg("resp timeout")
					} else {
						lg("Get error: %v", err)
					}
				} else {
					lg("ping")
				}

				time.Sleep(time.Duration(5) * time.Second)
			}
		}()
	} else {
		lg("slave")
		//slave
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			lg("pong")
			w.Write([]byte("Hello, world"))
		})
	}

	lg("ifpass")

	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		lg("healthcheck")
		w.WriteHeader(http.StatusOK)
	})

	// fot health check
	http.ListenAndServe(":8080", nil)
}
