package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type DbConfig struct {
	PgHost    string `env:"PG_HOST,notEmpty"`
	PgPort    uint16 `env:"PG_PORT,notEmpty"`
	PgUser    string `env:"PG_USER,notEmpty"`
	PgPass    string `env:"PG_PASS,notEmpty" logging:"mask"`
	PgDbName  string `env:"PG_DB,notEmpty"`
	PgMaxConn uint16 `env:"PG_MAX_CONN" envDefault:"20"`
}
type Config struct {
	AppName string `env:"APP_NAME,notEmpty" envDefault:"semen-platform"`

	PublicApiHost string `env:"PUBLIC_API_HOST"`
	PublicApiPort int    `env:"PUBLIC_API_PORT,notEmpty" envDefault:"8080"`

	JWTSecret string `env:"JWT_SECRET,notEmpty"`

	RedisAddr            string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	KafkaBroker          string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	NotificationGRPCAddr string `env:"NOTIFICATION_GRPC_ADDR" envDefault:"localhost:50051"`

	PG DbConfig
}

func Load(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("Config.Load(): %w", err)
	}

	return nil
}
func (cfg *Config) PublicApiAddr() string {
	return fmt.Sprintf("%s:%d", cfg.PublicApiHost, cfg.PublicApiPort)
}
func (cfg *DbConfig) PgConnDsn() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s",
		cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDbName,
	)
}
