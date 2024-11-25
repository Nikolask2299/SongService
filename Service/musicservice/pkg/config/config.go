package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigPostgres struct {
	User string
	Password string
	DBName string
	Host string
	Port string
}

type ServerConfig struct {
	Host string
	Port string
}

type APIConfig struct {
	Server ServerConfig
}

type ConfigMigrator struct {
	MigrationsPath string
	MigrationsTable string
}


func ReturnedDatabase() (ConfigPostgres, error) {
	err := godotenv.Load(Dir("config.env"))
	if err!= nil {
        return ConfigPostgres{}, err
    }

	configPostgres := ConfigPostgres{
		User:     getEnv("POSTGRES_USER", ""),
        Password: getEnv("POSTGRES_PASSWORD", ""),
        DBName:   getEnv("POSTGRES_DB", ""),
        Host:     getEnv("POSTGRES_HOST", ""),
        Port:     getEnv("POSTGRES_PORT", ""),
	}

	return configPostgres, nil
}

func InitConfigMigration() (ConfigPostgres, ConfigMigrator, error) {
	err := godotenv.Load(Dir("config.env"))
	if err!= nil {
        return ConfigPostgres{}, ConfigMigrator{}, err
    }

	configPostgres := ConfigPostgres{
		User:     getEnv("POSTGRES_USER", ""),
        Password: getEnv("POSTGRES_PASSWORD", ""),
        DBName:   getEnv("POSTGRES_DB", ""),
        Host:     getEnv("POSTGRES_HOST", ""),
        Port:     getEnv("POSTGRES_PORT", ""),
	}

	configMigrator := ConfigMigrator{
        MigrationsPath: getEnv("MIGRATIONS_PATH", ""),
        MigrationsTable: getEnv("MIGRATIONS_TABLE", ""),
    }

	return configPostgres, configMigrator, nil
}

func RetuneServerConfig() (ServerConfig, APIConfig, error) {
	err := godotenv.Load(Dir("config.env"))
    if err!= nil {
        return ServerConfig{}, APIConfig{}, err
    }
    
    serverConfig := ServerConfig{
        Host: getEnv("SERVER_HOST", ""),
        Port: getEnv("SERVER_PORT", ""),
    }

    apiConfig := APIConfig{
        Server: ServerConfig{
			Host: getEnv("API_HOST", ""),
            Port: getEnv("API_PORT", ""),
        },
    }
	
    return serverConfig, apiConfig, nil
}

func getEnv(key string, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
		return value
    }

    return defaultVal
}


func Dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}    
    currentDir = strings.Replace(currentDir, filepath.Join("Service","musicservice","cmd", "migration"), "", -1)
	currentDir = strings.Replace(currentDir, filepath.Join("Service","musicservice","cmd"), "", -1)
    return filepath.Join(currentDir, envFile)
}