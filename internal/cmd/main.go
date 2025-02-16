package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/garaevmir/avitocoinstore/internal/handler"
	"github.com/garaevmir/avitocoinstore/internal/middleware"
	"github.com/garaevmir/avitocoinstore/internal/repository"
	"github.com/garaevmir/avitocoinstore/internal/service"
)

func main() {
	runtime.GOMAXPROCS((runtime.NumCPU()))
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	transactionRepo := repository.NewTransactionRepository(pool)
	inventoryRepo := repository.NewInventoryRepository(pool)
	shopService := service.NewShopService(userRepo, transactionRepo, inventoryRepo)

	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(1000)))
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	authHandler := handler.NewAuthHandler(userRepo, os.Getenv("JWT_SECRET"))
	coinHandler := handler.NewCoinHandler(transactionRepo, userRepo)
	infoHandler := handler.NewInfoHandler(userRepo, inventoryRepo, transactionRepo)
	shopHandler := handler.NewShopHandler(shopService)

	e.POST("/api/auth", authHandler.Login)

	api := e.Group("/api")
	api.Use(middleware.JWTAuth(os.Getenv("JWT_SECRET")))
	api.GET("/info", infoHandler.GetUserInfo)
	api.POST("/sendCoin", coinHandler.SendCoins)
	api.GET("/buy/:item", shopHandler.BuyItem)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	s := &http.Server{
		Addr: ":8080",
	}

	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("HTTP server shutdown error:", err)
	}

	select {
	case <-ctx.Done():
		log.Printf("timeout of 5 seconds.\n")
	}
	log.Printf("Server exiting\n")

}
