package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wbl0/internal/cache"
	"wbl0/internal/config"
	"wbl0/internal/order"
	"wbl0/internal/stan"
	"wbl0/internal/storage"
	"wbl0/internal/transport"
)

var db *sqlx.DB

func main() {
	// Загрузка конфигурации
	cfg := config.New()

	// Инициализация строки подключения
	connectionString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.DB.User, cfg.DB.Password, cfg.DB.DatabaseName, cfg.DB.Host, cfg.DB.Port,
	)

	// Инициализация хранилища и кэша
	storage := storage.New(connectionString)
	cache := cache.New()
	log.Output(1, "sd")

	// Создание сервиса заказов
	orderSvc := order.New(storage, cache)

	// Создание и запуск сервиса STAN для работы с NATS Streaming
	stanService := stan.New(cfg.Stan, orderSvc)
	// Запуск http
	app := transport.New(cfg, orderSvc, stanService)
	go func() {
		if err := app.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error on start web server: %s\n", err)
		}
	}()

	//Завершение работы
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}
