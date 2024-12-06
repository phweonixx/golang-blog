package main

import (
	"blogAPI/internal/api"
	"blogAPI/internal/config"
	"blogAPI/internal/database"
	"blogAPI/internal/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	http.Handle("/", api.Router)

	cfg := config.New()

	// Функція з ініціалізіцією бази даних
	db := database.New()
	if db.DBGorm != nil {
		log.Println("DB initialized and ready to use.")
	}

	// Маршрути
	routes.SetupRoutes(api.Router)

	// Для правильної зупинки серверу
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    cfg.Config.Port,
		Handler: api.Router,
	}

	go func() {
		log.Println("Server started at", cfg.Config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
