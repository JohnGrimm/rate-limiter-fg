package cmd

import (
	"log"
	"net/http"

	"github.com/JohnGrimm/rate-limiter-fg/internal/config"
	"github.com/JohnGrimm/rate-limiter-fg/internal/handler"
	"github.com/JohnGrimm/rate-limiter-fg/internal/middleware"
	"github.com/JohnGrimm/rate-limiter-fg/internal/services"
	"github.com/JohnGrimm/rate-limiter-fg/internal/storage"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func Main() {

	cfg := config.GetConfig()

	// Usar as configurações
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	storageRedis := storage.NewRedisStorage(rdb)

	r := mux.NewRouter()
	limiter := *services.NewLimiter(storageRedis,
		config.GetConfig().DefaultRateLimit,
		config.GetConfig().DefaultExpiry,
		config.GetConfig().DefaultTimeBlocked)

	limiter.ProcessKeysFromFile()
	middleware := middleware.NewRateLimiterMiddleware(limiter)

	r.Use(middleware.Middleware)
	r.HandleFunc("/", handler.HandlerHello).Methods("GET")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
}
