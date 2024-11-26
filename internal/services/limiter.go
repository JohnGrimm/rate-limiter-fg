package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/JohnGrimm/rate-limiter-fg/internal/storage"
)

type Limiter struct {
	storage     storage.Storage
	fillRate    int // Tokens por segundo
	capacity    int // Capacidade máxima do balde
	timeBlocked int // Tempo em minutos que o ip ficarábloqueado
	keys        []ConfigKey
}

type ConfigKey struct {
	Key          string `json:"token"`
	TimeToExpire int    `json:"expiresIn"`
	QtdRequests  int    `json:"qtdRequests"`
}

const (
	Prefix_Key_TimeBlocked = "TB-"
	File_Config_Key        = "token-list.json"
)

func NewLimiter(storage storage.Storage, fillRate, capacity int, timeBlocked int) *Limiter {
	return &Limiter{storage: storage, fillRate: fillRate, capacity: capacity, timeBlocked: timeBlocked}
}

func (l *Limiter) Allow(key string, isIp bool) bool {
	ctx := context.Background()

	// Obter o número atual de tokens
	tokens, err := l.storage.Get(ctx, key)
	log.Printf("token consultado %s , tokens: %v", key, tokens)
	if err != nil {
		log.Printf("error to get key in redis, error: %v, config redis %v", err, l)
		if errors.As(err, &storage.ErrNotFound) && isIp {
			l.storage.Increment(ctx, key, l.fillRate)
			return l.Allow(key, isIp)
		}
		return false
	}

	// Verificar se há tokens suficientes
	if tokens <= 0 {
		return false
	}

	// Remover um token
	err = l.storage.Decrement(ctx, key)
	log.Println("foi feito o drecemento")
	if err != nil {
		// Lidar com erro
		return false
	}

	// Agendar a reposição de tokens
	go func() {
		timeBlocked := l.timeBlocked
		timeBlockedRedis, err := l.storage.Get(ctx, Prefix_Key_TimeBlocked+key)
		if err == nil {
			timeBlocked = timeBlockedRedis
		}
		log.Printf("time to blocked request, %v", timeBlocked)
		time.Sleep(time.Second * time.Duration(timeBlocked))
		l.storage.Increment(ctx, key, 1)
	}()

	return true
}

func (l *Limiter) ProcessKeysFromFile() {

	l.keys = l.loadKeyFromFile()

	log.Printf("loaded keys, %v", l.keys)
	for _, value := range l.keys {
		l.storage.Increment(context.Background(),
			Prefix_Key_TimeBlocked+value.Key,
			value.TimeToExpire)
		l.storage.Increment(context.Background(),
			value.Key,
			value.QtdRequests)
	}

}

func (l *Limiter) CompareKey(key string) bool {
	response := false
	for _, value := range l.keys {
		if value.Key == key {
			response = true
			break
		}
	}

	return response
}

func (l *Limiter) loadKeyFromFile() []ConfigKey {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error to get key from file token-list.json, %v", err)
	}

	filePath := filepath.Join(currentDir, File_Config_Key)

	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error to open config key, %v", err)
	}

	defer jsonFile.Close()

	var keys []ConfigKey
	if err := json.NewDecoder(jsonFile).Decode(&keys); err != nil {
		log.Fatalf("error to load keys, %v", err)
	}

	return keys
}
