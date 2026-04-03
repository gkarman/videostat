package config

import "fmt"

type Config struct {
	Env         string `env:"ENV" env-default:"local"`
	DB          DBConfig
	Logger      LoggerConfig
	ServerHttp  ServerHttpConfig
	ServerGRPC  ServerGRPCConfig
	RabbitMQ    RabbitMQConfig
	Mail        MailConfig
	TelegramBot TelegramBotConfig
}

type DBConfig struct {
	Host                         string `env:"DB_HOST"`
	Port                         int    `env:"DB_PORT"`
	User                         string `env:"DB_USER"`
	Password                     string `env:"DB_PASS"`
	Name                         string `env:"DB_NAME"`
	SSLMode                      string `env:"DB_SSLMODE"`
	MaxConnections               int32  `env:"DB_MAX_CONNECTIONS"`
	MinConnections               int32  `env:"DB_MIN_CONNECTIONS"`
	MaxConnectionLifeTimeMinutes int    `env:"DB_MAX_CONNECTION_LIFETIME_MINUTES"`
	MaxConnectionIdleTimeMinutes int    `env:"DB_MAX_CONNECTION_IDLE_TIME_MINUTES"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL"`
}

type ServerHttpConfig struct {
	Addr                string `env:"SERVER_HTTP_ADDR" env-default:":8080"`
	ReadTimeoutSeconds  int    `env:"SERVER_HTTP_READ_TIMEOUT_SECONDS" env-default:"10"`
	WriteTimeoutSeconds int    `env:"SERVER_HTTP_WRITE_TIMEOUT_SECONDS" env-default:"10"`
}

type ServerGRPCConfig struct {
	Addr string `env:"SERVER_GRPC_ADDR" env-default:"localhost:5051"`
}

type RabbitMQConfig struct {
	User           string `env:"RABBITMQ_USER" env-default:"guest"`
	Password       string `env:"RABBITMQ_PASS" env-default:"guest"`
	Host           string `env:"RABBITMQ_HOST" env-default:"localhost"`
	Port           string `env:"RABBITMQ_PORT" env-default:"5672"`
	Exchange       string `env:"RABBITMQ_EXCHANGE" env-default:"5672"`
	PoolSize       int    `env:"RABBITMQ_POOL_SIZE" env-default:"10"`
	ReconnectDelay int    `env:"RABBITMQ_RECONNECT_DELAY_IN_SECONDS" env-default:"3"`
}

type MailConfig struct {
	Driver   string `env:"MAIL_DRIVER" env-default:"smtp"`
	Host     string `env:"MAIL_HOST" env-default:"localhost"`
	Port     string `env:"MAIL_PORT" env-default:"2525"`
	User     string `env:"MAIL_USER" env-default:"user"`
	Password string `env:"MAIL_PASSWORD" env-default:"password"`
}

type TelegramBotConfig struct {
	Token string `env:"TELEGRAM_BOT_TOKEN" env-default:""`
	Debug bool `env:"TELEGRAM_BOT_DEBUG" env-default:"false""`
	Timeout int `env:"TELEGRAM_BOT_TIMEOUT" env-default:"60""`
}

