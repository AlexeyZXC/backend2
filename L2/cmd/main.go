package main

import (
	"context"
	"fmt"
	"k8s-go-app/config"
	"k8s-go-app/server"
	"k8s-go-app/version"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	launchMode := config.LaunchMode(os.Getenv("LAUNCH_MODE"))
	if len(launchMode) == 0 {
		launchMode = config.LocalEnv
	}
	log.Printf("LAUNCH MODE: %v", launchMode)
	cfg, err := config.Load(launchMode, "./config")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CONFIG: %+v", cfg)
	info := server.VersionInfo{
		Version: version.Version,
		Commit:  version.Commit,
		Build:   version.Build,
	}
	srv := server.New(info, cfg.Port)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err := srv.Serve(ctx)
		if err != nil {
			log.Println(fmt.Errorf("serve: %w", err))
			return
		}
	}()
	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-osSigChan
	log.Println("OS interrupting signal has received")
	cancel()
}
