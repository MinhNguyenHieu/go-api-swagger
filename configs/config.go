package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
	SMTP        SMTPConfig
	Auth        AuthConfig
	// RedisCfg    RedisConfig
	RateLimiter RateLimiterConfig
}

type SMTPConfig struct {
	Host        string
	Port        int
	User        string
	Pass        string
	SenderEmail string
}

type AuthConfig struct {
	Basic BasicAuthConfig
}

type BasicAuthConfig struct {
	User string
	Pass string
}

// type RedisConfig struct {
// 	Enabled  bool
// 	Addr     string
// 	Password string
// 	DB       int
// 	PoolSize int
// }

type RateLimiterConfig struct {
	Enabled bool
	RPS     float64
	Burst   int
	TTL     time.Duration
}

func LoadConfig() *Config {
	dbUser := getEnv("POSTGRES_USER", "user")
	dbPassword := getEnv("POSTGRES_PASSWORD", "password")
	dbName := getEnv("POSTGRES_DB", "mydatabase")
	dbHost := getEnv("POSTGRES_HOST", "db")
	dbPort := getEnv("POSTGRES_PORT", "5432")

	appPort := getEnv("APP_PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "secret_jwt_key")

	smtpHost := getEnv("SMTP_HOST", "")
	smtpPortStr := getEnv("SMTP_PORT", "0")
	smtpUser := getEnv("SMTP_USER", "")
	smtpPass := getEnv("SMTP_PASS", "")
	senderEmail := getEnv("SENDER_EMAIL", "")

	basicAuthUser := getEnv("BASIC_AUTH_USER", "admin")
	basicAuthPass := getEnv("BASIC_AUTH_PASS", "password")

	// redisEnabledStr := getEnv("REDIS_ENABLED", "false")
	// redisEnabled, err := strconv.ParseBool(redisEnabledStr)
	// if err != nil {
	// 	log.Printf("Warning: Invalid REDIS_ENABLED value, using false: %v", err)
	// 	redisEnabled = false
	// }
	// redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	// redisPassword := getEnv("REDIS_PASSWORD", "")
	// redisDBStr := getEnv("REDIS_DB", "0")
	// redisDB, err := strconv.Atoi(redisDBStr)
	// if err != nil {
	// 	log.Printf("Warning: Invalid REDIS_DB value, using 0: %v", err)
	// 	redisDB = 0
	// }
	// redisPoolSizeStr := getEnv("REDIS_POOL_SIZE", "10")
	// redisPoolSize, err := strconv.Atoi(redisPoolSizeStr)
	// if err != nil {
	// 	log.Printf("Warning: Invalid REDIS_POOL_SIZE value, using 10: %v", err)
	// 	redisPoolSize = 10
	// }

	rateLimiterEnabledStr := getEnv("RATE_LIMITER_ENABLED", "true")
	rateLimiterEnabled, err := strconv.ParseBool(rateLimiterEnabledStr)
	if err != nil {
		log.Printf("Warning: Invalid RATE_LIMITER_ENABLED value, using true: %v", err)
		rateLimiterEnabled = true
	}
	rateLimiterRPSStr := getEnv("RATE_LIMITER_RPS", "10")
	rateLimiterRPS, err := strconv.ParseFloat(rateLimiterRPSStr, 64)
	if err != nil {
		log.Printf("Warning: Invalid RATE_LIMITER_RPS value, using 10: %v", err)
		rateLimiterRPS = 10
	}
	rateLimiterBurstStr := getEnv("RATE_LIMITER_BURST", "10")
	rateLimiterBurst, err := strconv.Atoi(rateLimiterBurstStr)
	if err != nil {
		log.Printf("Warning: Invalid RATE_LIMITER_BURST value, using 10: %v", err)
		rateLimiterBurst = 10
	}
	rateLimiterTTLStr := getEnv("RATE_LIMITER_TTL", "1m")
	rateLimiterTTL, err := time.ParseDuration(rateLimiterTTLStr)
	if err != nil {
		log.Printf("Warning: Invalid RATE_LIMITER_TTL value, using 1m: %v", err)
		rateLimiterTTL = time.Minute
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("Warning: Invalid SMTP port, using 0: %v", err)
		smtpPort = 0
	}

	return &Config{
		AppPort: appPort,
		DatabaseURL: fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser, dbPassword, dbHost, dbPort, dbName),
		JWTSecret: jwtSecret,
		SMTP: SMTPConfig{
			Host:        smtpHost,
			Port:        smtpPort,
			User:        smtpUser,
			Pass:        smtpPass,
			SenderEmail: senderEmail,
		},
		Auth: AuthConfig{
			Basic: BasicAuthConfig{
				User: basicAuthUser,
				Pass: basicAuthPass,
			},
		},
		// RedisCfg: RedisConfig{
		// 	Enabled:  redisEnabled,
		// 	Addr:     redisAddr,
		// 	Password: redisPassword,
		// 	DB:       redisDB,
		// 	PoolSize: redisPoolSize,
		// },
		RateLimiter: RateLimiterConfig{
			Enabled: rateLimiterEnabled,
			RPS:     rateLimiterRPS,
			Burst:   rateLimiterBurst,
			TTL:     rateLimiterTTL,
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Printf("Warning: Environment variable '%s' not set, using default value: '%s'", key, fallback)
	return fallback
}
