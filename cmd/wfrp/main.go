package main

import (
	"fmt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/golangjwt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/mockemail"
	"github.com/jmilosze/wfrp-hammergen-go/internal/http"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	jwtService := golangjwt.NewHmacService(cfg.JwtConfig.HmacSecret, cfg.JwtConfig.AccessExpiryTime, cfg.JwtConfig.ResetExpiryTime)
	emailService := mockemail.NewEmailService()
	userService := memdb.NewUserService(cfg.MemDbUserService, cfg.MemDbUserService.SeedUsers, emailService, jwtService)

	router := gin.NewRouter()
	gin.RegisterUserRoutes(router, userService, jwtService)
	gin.RegisterAuthRoutes(router, userService, jwtService)

	server := http.NewServer(cfg.ServerConfig, router)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	server.Start()

	<-done

	if err := server.Stop(); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	return nil
}
