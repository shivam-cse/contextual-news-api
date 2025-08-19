package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shivam-cse/contextual-news-api/pkg/startup"
	"github.com/shivam-cse/contextual-news-api/pkg/logger"
	"github.com/shivam-cse/contextual-news-api/internal/dbInterface"
	services "github.com/shivam-cse/contextual-news-api/internal/services"
	v1Handlers "github.com/shivam-cse/contextual-news-api/internal/handlers/v1"
)

func Run() {

	// Load configuration first
	config, err := startup.LoadConfig()
    if err != nil {
        panic(err)
    }
	// Print all configuration settings
	// Just for debugging purposes, can be removed in production
	fmt.Printf("\n=== Configuration Settings ===\n")
	fmt.Printf("ServerAddress: %s\n", config.ServerAddress)
	fmt.Printf("ServerPort: %s\n", config.ServerPort)
	fmt.Printf("MongoConnectionString: %s\n", config.MongoConnectionString)
	fmt.Printf("MongoDatabase: %s\n", config.MongoDatabase)
	fmt.Printf("LLMToken: %s\n", config.LLMToken)
	fmt.Printf("LLMEndpoint: %s\n", config.LLMEndpoint)
	fmt.Printf("LLMModel: %s\n", config.LLMModel)
	fmt.Printf("\n===============================\n")

	// Create the logger
	logger := logger.New()
	logger.Info("Logger initialized")

	// Connect to MongoDB
	mongoClient, err := startup.ConnectMongoDB(config.MongoConnectionString)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		panic(err)
	}
	defer startup.Close(mongoClient)
	logger.Info("Successfully connected to MongoDB")

	// Create the news database
    database := mongoClient.Database(config.MongoDatabase)

	// Create the indexes to improve query performance
	err = startup.CreateIndexOnNewsColl(database)
	if err != nil {
		logger.Error("Failed to create indexes", "error", err)
		panic(err)
	}
	logger.Info("Indexes on news collection created successfully")

	// Create the news database interface
	newsDbInterface := dbInterface.NewNewsDbInterface(database, logger)

	// Create the LLM service
	llmService, err := services.NewLLMOpenRouterService(config.LLMToken, config.LLMEndpoint, config.LLMModel, logger)
	if err != nil {
		logger.Error("Failed to create LLM service", "error", err)
		panic(err)
	}
	logger.Info("LLM service created successfully", "model", config.LLMModel)

	// Create the news service
	newsService := services.NewNewsService(newsDbInterface, logger, llmService)

	// Create the news handler
	v1NewsHandler := v1Handlers.NewNewsHandler(newsService, logger)

	// Set up the router
	router := gin.Default()

	// Register the routes for v1
	v1Handlers.RegisterRoutes(router, v1NewsHandler)
	
	// Register the routes for v2
	// v2Handlers.RegisterRoutes(router, v2NewsHandler)

	err = router.Run(config.ServerAddress + ":" + config.ServerPort)
	
	if err != nil {
		logger.Error("Error starting server", "error", err)
		panic(err)
	}

	logger.Info("Server running on", "address", config.ServerAddress, "port", config.ServerPort)
}