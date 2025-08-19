package services

import (
	"fmt"
	"log/slog"
	"context"
	"encoding/json"
	"strings"
	"github.com/eduardolat/openroutergo"
	"github.com/shivam-cse/contextual-news-api/internal/models/newsArticle"
)

type LLMOpenRouterService struct {
	client   *openroutergo.Client
	Logger   *slog.Logger
	Token    string
	Endpoint string
	LLMModel string
}

func NewLLMOpenRouterService(
	token string,
	endpoint string,
	llmModel string,
	logger *slog.Logger,
) (*LLMOpenRouterService, error) {
	client, err := openroutergo.NewClient().
		WithAPIKey(token).
		WithBaseURL(endpoint).
		Create()
	if err != nil {
		return nil, err
	}
	return &LLMOpenRouterService{
		client:  client,
		Logger:  logger,
		Token:   token,
		Endpoint: endpoint,
		LLMModel: llmModel,
	}, nil
}

func (llmService *LLMOpenRouterService) GenerateSummary(
	ctx context.Context,
	systemMessage string,
	userMessage string,
) (string, error) {

	_, resp, err := llmService.client.
		NewChatCompletion().
		WithContext(ctx).
		WithModel(llmService.LLMModel).
		WithSystemMessage(systemMessage).
		WithUserMessage(userMessage).
		Execute()

	if err != nil {
		return "", err
	}

	llmService.Logger.Debug(fmt.Sprintf("Generated summary response: %v", resp.Choices[0].Message.Content))
	return resp.Choices[0].Message.Content, nil
}


func (llmService *LLMOpenRouterService) ExtractEntitiesAndIntent(
	ctx context.Context,
	systemMessage string,
	userMessage string,
) (newsArticle.LLMEntitiesAndIntentOutput, error) {

	_, resp, err := llmService.client.
		NewChatCompletion().
		WithContext(ctx).
		WithModel(llmService.LLMModel).
		WithSystemMessage(systemMessage).
		WithUserMessage(userMessage).
		Execute()

	if err != nil {
		return newsArticle.LLMEntitiesAndIntentOutput{}, err
	}
	llmService.Logger.Debug(fmt.Sprintf("Extracted entities and intent from user query response: %v", resp.Choices[0].Message.Content))

	var llmOutput newsArticle.LLMEntitiesAndIntentOutput

	// Extract JSON from markdown code blocks if present
	content := resp.Choices[0].Message.Content
	if strings.Contains(content, "```json") {
		start := strings.Index(content, "```json") + 7
		end := strings.LastIndex(content, "```")
		if start < end && end != -1 {
			content = strings.TrimSpace(content[start:end])
		}
	}
	
	if err := json.Unmarshal([]byte(content), &llmOutput); err != nil {
		return newsArticle.LLMEntitiesAndIntentOutput{}, err
	}

	return llmOutput, nil
}