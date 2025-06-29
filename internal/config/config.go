package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// Application settings
	App AppConfig `mapstructure:"app"`
	
	// Telegram Bot settings
	Telegram TelegramConfig `mapstructure:"telegram"`
	
	// Database settings
	Database DatabaseConfig `mapstructure:"database"`
	
	// Logging settings
	Logging LoggingConfig `mapstructure:"logging"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

type TelegramConfig struct {
	BotToken    string `mapstructure:"bot_token"`
	WebhookURL  string `mapstructure:"webhook_url"`
	UseWebhook  bool   `mapstructure:"use_webhook"`
	AllowedUser int64  `mapstructure:"allowed_user"` // Для персонального использования
}

type DatabaseConfig struct {
	URL             string `mapstructure:"url"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Name            string `mapstructure:"name"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	SSLMode         string `mapstructure:"ssl_mode"`
	MaxConnections  int    `mapstructure:"max_connections"`
	MaxIdleTime     int    `mapstructure:"max_idle_time"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	Filename string `mapstructure:"filename"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.subscription-tracker")

	// Set default values
	setDefaults()

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("") // No prefix, use direct env var names
	
	// Bind environment variables explicitly
	viper.BindEnv("telegram.bot_token", "TELEGRAM_BOT_TOKEN")
	viper.BindEnv("database.url", "DATABASE_URL")
	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.name", "DATABASE_NAME")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("app.environment", "APP_ENVIRONMENT")
	viper.BindEnv("app.debug", "APP_DEBUG")
	viper.BindEnv("logging.level", "LOGGING_LEVEL")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, continue with env vars and defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Build database URL if not provided
	if config.Database.URL == "" && config.Database.Host != "" {
		config.Database.URL = buildDatabaseURL(config.Database)
	}

	// Validate required fields
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "Subscription Tracker Bot")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", false)

	// Telegram defaults
	viper.SetDefault("telegram.use_webhook", false)
	viper.SetDefault("telegram.allowed_user", 0)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "sub_cos_counter")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_connections", 25)
	viper.SetDefault("database.max_idle_time", 15) // minutes
	viper.SetDefault("database.conn_max_lifetime", 60) // minutes

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
}

func buildDatabaseURL(db DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

func validateConfig(config *Config) error {
	if config.Telegram.BotToken == "" {
		return fmt.Errorf("telegram bot token is required")
	}

	if config.Database.URL == "" && config.Database.Host == "" {
		return fmt.Errorf("database configuration is required")
	}

	return nil
}

// Convenience methods
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

func (c *Config) GetDatabaseURL() string {
	return c.Database.URL
}

func (c *Config) GetBotToken() string {
	return c.Telegram.BotToken
}

func (c *Config) IsDebugMode() bool {
	return c.App.Debug || c.IsDevelopment()
}

func (c *Config) IsPersonalBot() bool {
	return c.Telegram.AllowedUser != 0
}