package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"pipo/version"
	"strings"
	"time"
)

// /home/user/sources/geekBrains/backend2/L2/pipo/main.version=1.0.c
// go build -ldflags="-X /home/user/sources/geekBrains/backend2/L2/pipo/main.version=1.0.c"

// go build -ldflags="-X main.version=1.0.d" //ok
// var version = "0.0.0" //ok
// var version string //ok

// const version = "0.0.0" // not work; version remains 0.0.0

// go build -ldflags="-X pipo/version.Version=1.0.d" //ok
// go build -ldflags="-X /home/user/sources/geekBrains/backend2/L2/pipo/pipo/version.Version=1.0.c" // does not work

func main() {
	var version = version.Version

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
				reqReader := strings.NewReader(fmt.Sprintf("ping(%v)", version))
				// resp, err := http.Get(fmt.Sprintf("http://%s:8080", slaveName))
				resp, err := http.Post(fmt.Sprintf("http://%s:8080", slaveName), "text/plain", reqReader)
				if err != nil {
					var urlErr url.Error
					if errors.Is(err, &urlErr) && urlErr.Timeout() {
						lg("resp timeout")
					} else {
						lg("Get error: %v", err)
					}
				} else {
					//response ok
					respBody, err := io.ReadAll(resp.Body)
					resp.Body.Close()

					if err != nil {
						lg("ReadAll error: %v", err)
					} else {
						lg("resp(%v)", string(respBody))
					}
				}

				time.Sleep(time.Duration(5) * time.Second)
			}
		}()
	} else {
		lg("slave")
		//slave
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			r.Body.Close()
			if err != nil {
				lg("ReadAll error: %v", err)
			} else {
				lg("got request(%v)", string(body))
				w.Write([]byte(fmt.Sprintf("pong(%v)", version)))
			}
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
