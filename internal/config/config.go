package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Discord  DiscordConfig  `json:"discord"`
	WhatsApp WhatsAppConfig `json:"whatsapp"`
	Jira     JiraConfig     `json:"jira"`
	Gemini   GeminiConfig   `json:"gemini"`
	Server   ServerConfig   `json:"server"`
}

type DiscordConfig struct {
	Token   string `json:"token"`
	GuildID string `json:"guild_id"`
}

type WhatsAppConfig struct {
	DatabaseURI string `json:"database_uri"`
	LogLevel    string `json:"log_level"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Database    string `json:"database"`
	User        string `json:"user"`
	Password    string `json:"password"`
	SSLMode     string `json:"ssl_mode"`
}

type JiraConfig struct {
	URL        string `json:"url"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ProjectKey string `json:"project_key"`
}

type GeminiConfig struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
}

type ServerConfig struct {
	Port        string `json:"port"`
	LogLevel    string `json:"log_level"`
	Environment string `json:"environment"`
}

func Load() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("error cargando .env: %w", err)
		}
	}

	cfg := &Config{}

	cfg.Discord = DiscordConfig{
		Token:   getEnv("DISCORD_TOKEN", ""),
		GuildID: getEnv("DISCORD_GUILD_ID", ""),
	}

	port, _ := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	cfg.WhatsApp = WhatsAppConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     port,
		Database: getEnv("POSTGRES_DB", "lisa_whatsmeow"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", ""),
		SSLMode:  getEnv("POSTGRES_SSL", "disable"),
		LogLevel: getEnv("LOG_LEVEL", "INFO"),
	}
	cfg.WhatsApp.DatabaseURI = buildPostgresURI(cfg.WhatsApp)

	cfg.Jira = JiraConfig{
		URL:        getEnv("JIRA_URL", ""),
		Email:      getEnv("JIRA_EMAIL", ""),
		Token:      getEnv("JIRA_TOKEN", ""),
		ProjectKey: getEnv("JIRA_PROJECT_KEY", "SUPPORT"),
	}

	cfg.Gemini = GeminiConfig{
		APIKey: getEnv("GEMINI_API_KEY", ""),
		Model:  getEnv("GEMINI_MODEL", "gemini-pro"),
	}

	cfg.Server = ServerConfig{
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "INFO"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuración inválida: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var missing []string

	if c.Discord.Token == "" {
		missing = append(missing, "DISCORD_TOKEN")
	}

	if c.WhatsApp.Password == "" {
		missing = append(missing, "POSTGRES_PASSWORD")
	}

	if c.Gemini.APIKey == "" {
		missing = append(missing, "GEMINI_API_KEY")
	}

	if c.Jira.URL == "" || c.Jira.Email == "" || c.Jira.Token == "" {
		fmt.Println("ADVERTENCIA: Configuración de Jira incompleta. Algunas funciones pueden no estar disponibles.")
	}

	if len(missing) > 0 {
		return fmt.Errorf("variables de entorno requeridas faltantes: %s", strings.Join(missing, ", "))
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

func (c *Config) GetServerAddress() string {
	return ":" + c.Server.Port
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildPostgresURI(cfg WhatsAppConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SSLMode,
	)
}
