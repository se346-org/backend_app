package configuration

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server        *ServerConfig        `yaml:"server,omitempty"`
	Postgres      *PostgresConfig      `yaml:"postgres,omitempty"`
	Redis         *RedisConfig         `yaml:"redis,omitempty"`
	WebSocket     *WebSocketConfig     `yaml:"websocket,omitempty"`
	JWT           *JWTConfig           `yaml:"jwt,omitempty"`
	Nats          *NatsConfig          `yaml:"nats,omitempty"`
	Observability *ObservabilityConfig `yaml:"observability,omitempty"`
	Minio         *MinioConfig         `yaml:"minio,omitempty"`
}

type ServerConfig struct {
	Port   int    `yaml:"port,omitempty"`
	Origin string `yaml:"origin,omitempty"`
}

type PostgresConfig struct {
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database string `yaml:"database,omitempty"`
	SSLMode  string `yaml:"sslmode,omitempty"`
}

type RedisConfig struct {
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database int    `yaml:"database,omitempty"`
}

type WebSocketConfig struct {
	Origin string `yaml:"origin,omitempty"`
}

type JWTConfig struct {
	SecretKey  string `yaml:"secret_key,omitempty"`
	Issuer     string `yaml:"issuer,omitempty"`
	Expiration int    `yaml:"expiration,omitempty"` // in seconds
}

type NatsConfig struct {
	Address string `yaml:"address,omitempty"`
}

type ObservabilityConfig struct {
	TracingEnabled bool   `yaml:"tracing_enabled,omitempty"`
	JaegerEndpoint string `yaml:"jaeger_endpoint,omitempty"`
	JaegerService  string `yaml:"jaeger_service,omitempty"`
}

type MinioConfig struct {
	Endpoint       string `yaml:"endpoint,omitempty"`
	AccessKey      string `yaml:"access_key,omitempty"`
	SecretKey      string `yaml:"secret_key,omitempty"`
	Token          string `yaml:"token,omitempty"`
	UseSSL         bool   `yaml:"use_ssl,omitempty"`
	PublicEndpoint string `yaml:"public_endpoint,omitempty"`
}

var ConfigInstance *Config

func LoadConfig(configFilePath string) error {
	// Read the yaml file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	// Set the global config instance
	ConfigInstance = &config
	return nil
}
