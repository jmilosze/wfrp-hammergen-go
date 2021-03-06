package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/gin"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/golangjwt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/mailjet"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/mockcaptcha"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/mongodb"
	"github.com/jmilosze/wfrp-hammergen-go/internal/http"
	"github.com/jmilosze/wfrp-hammergen-go/internal/services"
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

	cfg := config.NewConfig()

	val := validator.New()
	jwtService := golangjwt.NewHmacService(cfg.Jwt.HmacSecret, cfg.Jwt.AccessExpiry, cfg.Jwt.ResetExpiry)
	emailService := mailjet.NewEmailService(cfg.Email.FromAddress, cfg.Email.PublicApiKey, cfg.Email.PrivateApiKey)
	captchaService := mockcaptcha.NewCaptchaService()

	mongoDbService := mongodb.NewDbService(cfg.MongoDb.Uri, cfg.MongoDb.DbName)
	defer mongoDbService.Disconnect()
	userDbService := mongodb.NewUserDbService(mongoDbService, cfg.MongoDb.UserCollection, cfg.MongoDb.CreateIndexes)

	userService := services.NewUserService(cfg.UserService, userDbService, emailService, jwtService, val)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.RequestTimeout)
	defer cancel()

	if cfg.UserService.CreateMockUsers {
		mockUsers := config.NewMockUsers()
		userService.SeedUsers(ctx, mockUsers)
	}

	router := gin.NewRouter(cfg.Server.RequestTimeout)
	gin.RegisterUserRoutes(router, userService, jwtService, captchaService)
	gin.RegisterAuthRoutes(router, userService, jwtService)

	server := http.NewServer(cfg.Server, router)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	server.Start()
	<-done
	server.Stop()

	return nil
}
