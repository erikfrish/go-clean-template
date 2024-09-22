package config

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
)

// Конфиг загружается из файла, переменных окружения и struct-tag'ов "env-default" с приоритетом:
// 1. Переменные окружения
// 2. Файл конфигурации
// 3. Значения по умолчанию в структуре

type Config struct {
	AppName      string     `yaml:"app-name"    json:"app_name"     env:"app_name"`
	AppVersion   string     `yaml:"app-version" json:"app_version"  env:"app_version"`
	PromPrefix   string     `yaml:"prom-prefix" json:"prom_prefix"`
	Env          string     `yaml:"env"         json:"env"          env:"env"`
	InstanceID   uuid.UUID  `yaml:"instance-id" json:"instance_id"`
	Logger       Logger     `yaml:"logger"      json:"logger"`
	ConfigString string     `yaml:"-"           json:"-"`
	DB           DB         `yaml:"db"          json:"db"`
	HTTP         HTTP       `yaml:"http"        json:"http"`
	Schedules    Schedules  `yaml:"schedules"   json:"schedules"`
	HTTPClient   HTTPClient `yaml:"http-client" json:"http_client"`
	API          API        `yaml:"api"         json:"api"`
}

type HTTP struct {
	Port         string        `yaml:"port"          json:"port"          env:"http_server_port"`
	ReadTimeout  time.Duration `yaml:"read-timeout"  json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout" json:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle-timeout"  json:"idle_timeout"`
}

type DB struct {
	Enabled         bool          `yaml:"enabled"           json:"enabled"   env:"db_enabled"`
	Host            string        `yaml:"host"              json:"host"      env:"db_host"`
	Port            string        `yaml:"port"              json:"port"      env:"db_port"`
	Database        string        `yaml:"database"          json:"database"  env:"db_database"`
	Schema          string        `yaml:"schema"            json:"schema"    env:"db_schema"`
	Username        string        `yaml:"username"          json:"-"         env:"db_username"`
	Password        string        `yaml:"password"          json:"-"         env:"db_password"`
	Scheme          string        `yaml:"scheme"            json:"scheme"`
	Driver          string        `yaml:"driver"            json:"driver"`
	FailoverHost    string        `yaml:"failover-host"     json:"failover_host"`
	MaxIdleConns    int           `yaml:"max-idle-conns"    json:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max-open-conns"    json:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn-max-lifetime" json:"conn_max_lifetime"`
	SSLMode         bool          `yaml:"ssl-mode"          json:"ssl_mode"`
}

type Schedules struct {
	Persist string `yaml:"persist" json:"persist" env:"persist-schedule"`
}

type HTTPClient struct {
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
}

type API struct {
	URL  string `yaml:"url"  json:"url"  env:"api_url"`
	Path string `yaml:"path" json:"path" env:"path"`
}

func (c *Config) String() string {
	return c.ConfigString
}

type Logger struct {
	LoggerTelegram LoggerTelegram `yaml:"logger-telegram" json:"telegram"`
	LoggerStd      LoggerStd      `yaml:"logger-std"      json:"std"`
	LoggerSlog     LoggerSlog     `yaml:"logger-slog"     json:"slog"`
}

type LoggerTelegram struct {
	Enabled      bool   `yaml:"enabled"        json:"enabled" env:"telegram_enabled"`
	Level        string `yaml:"level"          json:"level"   env:"telegram_level"`
	TargetChatID int64  `yaml:"target-chat-id" json:"chat_id" env:"telegram_chat_id"`
	BotAPIToken  string `yaml:"bot-api-token"  json:"-"       env:"bot_api_token"`
}

type LoggerStd struct {
	Enabled bool   `yaml:"enabled"  json:"enabled" env:"std_enabled"`
	Level   string `yaml:"level"    json:"level"   env:"std_level"`
	LogFile string `yaml:"log-file" json:"file"`
	Stdout  bool   `yaml:"stdout"   json:"stdout"`
}

type LoggerSlog struct {
	Enabled bool   `yaml:"enabled" json:"enabled" env:"slog_enabled"`
	Level   string `yaml:"level"   json:"level"   env:"slog_level"`
	JSON    bool   `yaml:"json"    json:"json"    env:"slog_json"`
}

func Load(fileName string) (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig(fileName, &cfg)
	if err != nil {
		return nil, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		out = []byte("config marshal error: " + err.Error())
	}
	cfg.ConfigString = string(out)

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("cleanenv.ReadEnv: %w", err)
	}

	return &cfg, nil
}
