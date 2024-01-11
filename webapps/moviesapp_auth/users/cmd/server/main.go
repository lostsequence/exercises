package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"movies-auth/users/internal/api/handlers"
	"movies-auth/users/internal/api/middlewares"
	"movies-auth/users/internal/config"
	"movies-auth/users/internal/services"
	"movies-auth/users/internal/storage/db"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
		return
	}

	cfgPath := os.Getenv("CONFIG_PATH")
	cfgName := os.Getenv("CONFIG_NAME")

	viper.AddConfigPath(cfgPath)
	viper.SetConfigName(cfgName)

	err = viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}

	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Println(err)
		return
	}

	dbCon, err := sql.Open("pgx", cfg.DBConfig.ConnectionString())
	if err != nil {
		log.Println(err)
		return
	}
	err = dbCon.Ping()
	if err != nil {
		log.Printf("failed to connect to db %v", err)
		return
	}

	dbStorage := db.NewDbStorage(dbCon)
	usersService := services.NewUsersService(dbStorage)
	sessionsService := services.NewSessionService(dbStorage)
	usersHandler := handlers.NewUsersHandler(usersService, sessionsService)

	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.Auth(sessionsService))
		r.Post("/register", usersHandler.Register)
		r.Post("/login", usersHandler.Login)
		r.Post("/logout", usersHandler.Logout)
		r.Get("/list", usersHandler.List)
		r.Get("/sessions/{key}", usersHandler.Session)
	})

	addr := fmt.Sprintf("%s:%d", cfg.ServerConfig.Host, cfg.ServerConfig.Port)

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)

	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}
	log.Println("starting server...")
	go func() {
		err = srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server stopped")
			return
		}

		log.Printf("unexpected server error: %s", err)
	}()
	log.Printf("server started on: %s", addr)

	<-ctx.Done()
	stop()

	log.Println("stopping database...")
	dbCon.Close()
	log.Println("database stopped")

	tCtx, tCancel := context.WithTimeout(ctx, time.Second*30)
	defer tCancel()
	err = srv.Shutdown(tCtx)
	if err != nil {
		log.Printf("server shutdown error: %s", err)
	}
}
