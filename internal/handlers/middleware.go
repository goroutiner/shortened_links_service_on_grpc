package handlers

import (
	"context"
	"log"
	"shortened_links_service_on_grpc/internal/config"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	visitors = make(map[string]*visitor) // visitors словарь для связи ip -> visitor
	mu       sync.Mutex
)

// visitor внутренняя структура для хранения лимитера и времени последнего запроса
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// СleanupVisitors очищает словарь visitors через каждый временной интервал,
// если пользователь не активен (временные параметры задаются в congif/setting.go)
func СleanupVisitors() {
	ticker := time.NewTicker(config.CleanupInterval)
	for range ticker.C {
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > config.InactivityLimit {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// getVisitor записывает в словарь visitors лимитеры для заданного ip
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if exists {
		// Обновляем время последнего запроса пользователя
		v.lastSeen = time.Now()
		return v.limiter
	}

	v = &visitor{
		limiter:  rate.NewLimiter(config.RateLimit, config.BufferLimit),
		lastSeen: time.Now(),
	}
	visitors[ip] = v
	return v.limiter
}

// LimiterMiddleware проверяет не превышен ли RPS по IP пользователя
func RateLimitInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			log.Println("Could not get peer info")
			return nil, status.Errorf(codes.Canceled, "Could not determine client IP")
		}
		
		ip := p.Addr.String()
		limiter := getVisitor(ip)
		if !limiter.Allow() {
			log.Printf("Too Many Requests for the user: %s\n", ip)
			return nil, status.Errorf(codes.ResourceExhausted, "Too Many Requests for the user: %s", ip)
		}

		return handler(ctx, req)
	}
}
