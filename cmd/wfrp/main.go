package main

import (
	"fmt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/golangjwt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/mongodb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cfg, err := config.NewDefault()
	if err != nil {
		return fmt.Errorf("getting service config from environment: %w", err)
	}

	jwtService := golangjwt.NewHmacService("some secret", 60*time.Minute)
	userService := mongodb.NewUserService()

	router := gin.NewRouter()
	gin.RegisterUserRoutes(router, userService, jwtService)
	gin.RegisterAuthRoutes(router, userService, jwtService)

	server := http.NewServer(cfg.APIServer, router)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	server.Start()

	<-done

	if err := server.Stop(); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	return nil
}