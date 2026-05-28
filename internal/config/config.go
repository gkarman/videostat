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
	Apify       ApifyConfig
	Assemblyai  AssemblyaiConfig
	Anthropic   AnthropicConfig
	OpenRouter  OpenRouterConfig
	OpenAI      OpenAIConfig
	LLM         LLMConfig
	S3          S3Config
	HeyGen      HeyGenConfig
	Kling       KlingConfig
	Shotstack   ShotstackConfig
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
	Token   string `env:"TELEGRAM_BOT_TOKEN" env-default:""`
	ChatID  int64  `env:"TELEGRAM_BOT_CHAT_ID" env-default:"0"`
	Debug   bool   `env:"TELEGRAM_BOT_DEBUG" env-default:"false"`
	Timeout int    `env:"TELEGRAM_BOT_TIMEOUT" env-default:"60"`
}

type ApifyConfig struct {
	Token string `env:"APIFY_TOKEN" env-default:""`
	Host  string `env:"APIFY_HOST" env-default:""`

	YoutubeMaxVideos int `env:"APIFY_YOUTUBE_MAX_VIDEOS" env-default:"10"`
	YoutubeDays      int `env:"APIFY_YOUTUBE_DAYS" env-default:"10"`

	TiktokMaxVideos int `env:"APIFY_TIKTOK_MAX_VIDEOS" env-default:"10"`
	TiktokDays      int `env:"APIFY_TIKTOK_DAYS" env-default:"10"`

	InstagramMaxVideos int `env:"APIFY_INSTAGRAM_MAX_VIDEOS" env-default:"10"`
	InstagramDays      int `env:"APIFY_INSTAGRAM_DAYS" env-default:"10"`
}

type AssemblyaiConfig struct {
	Token string `env:"ASSEMBLYAI_TOKEN" env-default:""`
	Host  string `env:"ASSEMBLYAI_HOST" env-default:""`
}

type AnthropicConfig struct {
	Token string `env:"ANTHROPIC_TOKEN" env-default:""`
	Model string `env:"ANTHROPIC_MODEL" env-default:"claude-sonnet-4-6"`
}

type OpenRouterConfig struct {
	Token string `env:"OPENROUTER_TOKEN" env-default:""`
	Model string `env:"OPENROUTER_MODEL" env-default:"meta-llama/llama-3.3-70b-instruct:free"`
}

type OpenAIConfig struct {
	Token string `env:"OPENAI_TOKEN" env-default:""`
	Model string `env:"OPENAI_MODEL" env-default:"gpt-4o-mini"`
}

// LLMConfig selects which LLM provider is active.
// LLM_PROVIDER: "anthropic" | "openrouter" | "openai"
type LLMConfig struct {
	Provider string `env:"LLM_PROVIDER" env-default:"anthropic"`
}

type HeyGenConfig struct {
	APIKey   string `env:"HEYGEN_API_KEY" env-default:""`
	AvatarID string `env:"HEYGEN_AVATAR_ID" env-default:""`
	VoiceID  string `env:"HEYGEN_VOICE_ID" env-default:""`
}

type ShotstackConfig struct {
	APIKey  string `env:"SHOTSTACK_API_KEY" env-default:""`
	BaseURL string `env:"SHOTSTACK_BASE_URL" env-default:"https://api.shotstack.io/edit/stage"`
}

type KlingConfig struct {
	AccessKeyID string `env:"KLING_ACCESS_KEY_ID" env-default:""`
	SecretKey   string `env:"KLING_SECRET_KEY" env-default:""`
	Model       string `env:"KLING_MODEL" env-default:"kling-v1"`
}

type S3Config struct {
	Endpoint  string `env:"S3_ENDPOINT" env-default:"http://localhost:9000"`
	PublicURL string `env:"S3_PUBLIC_URL" env-default:""`  // ngrok locally, same as Endpoint in prod
	AccessKey string `env:"S3_ACCESS_KEY" env-default:"minioadmin"`
	SecretKey string `env:"S3_SECRET_KEY" env-default:"minioadmin"`
	Bucket    string `env:"S3_BUCKET" env-default:"videostat"`
	Region    string `env:"S3_REGION" env-default:"us-east-1"`
}
