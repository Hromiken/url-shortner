package internal

import (
	"fmt"
	"net/http"
	config2 "shortner/config"
	"shortner/internal/handler"
	"shortner/internal/pkg/infrastructure"
	"shortner/internal/repositories/postgres"
	"shortner/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func Run() {
	// config
	cfg, err := config2.NewAppConfig()
	if err != nil {
		panic(err)
	}
	// logger
	zlog.Init()
	err = zlog.SetLevel(cfg.Logger.LogLevel)
	if err != nil {
		panic(err)
	}
	zlog.Logger.Info()

	//Postgres
	db, err := infrastructure.NewPostgres(cfg.Postgres)
	if err != nil {
		zlog.Logger.Fatal()
	}
	defer db.Master.Close()

	// repositories
	store := postgres.NewURLRepo(db)

	cache := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	// service
	svc := service.NewShorterService(store)
	svc.Cache = cache

	// handler
	h := handler.NewHandler(svc)

	// router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/shorten", h.SaveURL)
	router.Get("/s/", h.Redirect)
	router.Get("/analytics/{alias}", h.GetAnalytics)

	zlog.Logger.Info()
	err = http.ListenAndServe(cfg.Server.Addr, router)
	fmt.Printf("Starting server on addr=%q\n", cfg.Server.Addr)
	if err != nil {
		panic(err)
	}

}
