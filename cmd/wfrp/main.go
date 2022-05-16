package main

import (
	"fmt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/config"
	"github.com/jmilosze/wfrp-hammergen-go/internal/dependencies/mongodb"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cfg, err := config.NewFromEnv()
	if err != nil {
		return fmt.Errorf("getting service config from environment: %w", err)
	}

	//ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
	//defer cancel()
	//
	//val := validator.New()
	//jwtService := golangjwt.NewHmacService(cfg.JwtConfig.HmacSecret, cfg.JwtConfig.AccessExpiry, cfg.JwtConfig.ResetExpiry)
	//emailService := mockemail.NewEmailService(cfg.EmailConfig.FromAddress)
	//captchaService := mockcaptcha.NewCaptchaService()
	mongoDbService := mongodb.NewDbService(cfg.MongoDbConfig.Uri)
	defer mongoDbService.Disconnect()
	//userService := services.NewUserService(ctx, cfg.UserServiceConfig, userDbService, emailService, jwtService, val)
	//
	//router := gin.NewRouter(cfg.RequestTimeout)
	//gin.RegisterUserRoutes(router, userService, jwtService, captchaService)
	//gin.RegisterAuthRoutes(router, userService, jwtService)

	//server := http.NewServer(cfg.ServerConfig, router)
	//
	//done := make(chan os.Signal, 1)
	//signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	//
	//server.Start()
	//
	//<-done
	//
	//if err := server.Stop(); err != nil {
	//	log.Fatalf("Server Shutdown Failed:%+v", err)
	//}

	return nil
}
