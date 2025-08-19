package startup

import (
	"os"
	"strconv"
	"time"
	"github.com/joho/godotenv"
)

var ENV_DIR = "../../.env"

type Config struct {
	ServerAddress        	string
	ServerPort            	string
	MongoConnectionString 	string
	MongoDatabase         	string
	HandlerTimeoutDuration  time.Duration
	LLMToken               	string
	LLMEndpoint            	string
	LLMModel               	string
}

func LoadConfig(path ...string) (*Config, error) {
	if len(path) > 0 {
		ENV_DIR = path[0]
	}

	// Load environment variables from the .env file
	err := godotenv.Load(ENV_DIR)
	if err != nil {
		return nil, err
	}

	timeoutStr := getEnv("HANDLER_TIMEOUT_DURATION", "10")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 10 // default value
	}
	
	return &Config{
		ServerAddress:          getEnv("SERVER_ADDRESS", "localhost"),
		ServerPort:             getEnv("SERVER_PORT", "8080"),
		MongoConnectionString:  getEnv("MONGO_CONNECTION_STRING", "mongodb://localhost:27017"),
		MongoDatabase:          getEnv("MONGO_DATABASE", "contextual_news_api"),
		HandlerTimeoutDuration: time.Duration(timeout) * time.Second,
		LLMToken:               getEnv("LLM_TOKEN", ""),
		LLMEndpoint:            getEnv("LLM_ENDPOINT", ""),
		LLMModel:               getEnv("LLM_MODEL", "gpt-4o"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
