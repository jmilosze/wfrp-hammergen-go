package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/golangjwt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/memdb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/mockcaptcha"
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

	validate := validator.New()
	jwtService := golangjwt.NewHmacService(cfg.JwtConfig.HmacSecret, cfg.JwtConfig.AccessExpiryTime, cfg.JwtConfig.ResetExpiryTime)
	emailService := mockemail.NewEmailService(cfg.EmailConfig.FromAddress)
	captchaService := mockcaptcha.NewCaptchaService()
	userService := memdb.NewUserService(cfg.MemDbUserService, emailService, jwtService, validate)

	router := gin.NewRouter()
	gin.RegisterUserRoutes(router, userService, jwtService, captchaService)
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
