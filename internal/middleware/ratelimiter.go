package middleware

import (
	"net/http"
	"strings"

	"github.com/JohnGrimm/rate-limiter-fg/internal/services"
)

const HeaderAPIKey = "API_KEY"

type rateLimiterMiddleware struct {
	limiter services.Limiter
}

func NewRateLimiterMiddleware(limiter services.Limiter) *rateLimiterMiddleware {
	return &rateLimiterMiddleware{limiter: limiter}
}

func (m *rateLimiterMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter a chave para o rate limiting (IP ou token)
		key, isIp := m.getLimiterKey(r) // Função para extrair a chave do request

		// Verificar se a requisição está dentro do limite
		if !m.limiter.Allow(key, isIp) {
			http.Error(w, " you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		// Chamar o próximo handler
		next.ServeHTTP(w, r)
	})
}

// Função auxiliar para obter a chave para o rate limiting
func (m *rateLimiterMiddleware) getLimiterKey(r *http.Request) (string, bool) {
	// Implemente a lógica para extrair a chave, seja ela o IP ou o token
	if token := r.Header.Get(HeaderAPIKey); token != "" {
		// Extrair o token da autorização
		return token, false
	}
	return strings.Split(r.RemoteAddr, ":")[0], true // Usar o IP como chave
}
