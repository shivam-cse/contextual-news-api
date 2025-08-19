package services

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/shivam-cse/contextual-news-api/internal/dbInterface"
	"github.com/shivam-cse/contextual-news-api/internal/models/newsArticle"
	"github.com/shivam-cse/contextual-news-api/pkg/constants"
	"github.com/shivam-cse/contextual-news-api/pkg/utils"
)

type NewsService struct {
	DbInterface *dbInterface.NewsDbInterface
	Logger      *slog.Logger
	LLMService  *LLMOpenRouterService
}

func NewNewsService(
	dbInterface *dbInterface.NewsDbInterface,
	logger *slog.Logger,
	llmService *LLMOpenRouterService,
) *NewsService {
	return &NewsService{
		DbInterface: dbInterface,
		Logger:      logger,
		LLMService:  llmService,
	}
}

func (service *NewsService) ArticleSummaryHelper(
	ctx context.Context,
	articles []newsArticle.NewsArticleDBResponse,
) ([]newsArticle.NewsArticleDBResponse, error) {

	systemPrompt := constants.ARTICLE_NEWS_SUMMARY_SYSTEM_PROMPT

	// Call the LLM service for article summaries
	for i := range articles {
		article := articles[i]
		userPrompt := fmt.Sprintf(constants.ARTICLE_NEWS_SUMMARY_USER_PROMPT, article.Title, article.Description)

		summary, err := service.LLMService.GenerateSummary(ctx, systemPrompt, userPrompt)
		if err != nil {
			service.Logger.Error("Failed to get article summary", "error", err)
			return nil, err
		}
		articles[i].LLMSummary = summary
	}

	return articles, nil
}

func (service *NewsService) LatestNewsService(
	ctx context.Context,
	articleLimit int,
) ([]newsArticle.NewsArticleDBResponse, error) {
	service.Logger.Debug("'Service Layer': Fetching latest news articles...")

	articles, err := service.DbInterface.FindAllArticles(ctx, constants.NEWS, int64(articleLimit))
	if err != nil {
		service.Logger.Error("Failed to fetch latest news articles", "error", err)
		return nil, err
	}

	service.Logger.Info(fmt.Sprintf("Fetched %d latest news articles from database and creating summaries...", len(articles)))
	// Summarize the articles
	articles, err = service.ArticleSummaryHelper(ctx, articles)
	if err != nil {
		service.Logger.Error("Failed to summarize articles", "error", err)
		return nil, err
	}
	service.Logger.Info(fmt.Sprintf("Summarized %d latest news articles", len(articles)))

	return articles, nil
}

func (service *NewsService) CategoryNewsService(
	ctx context.Context,
	category string,
	articleLimit int,
) ([]newsArticle.NewsArticleDBResponse, error) {
	service.Logger.Debug("'Service Layer': Fetching news articles by category...")

	articles, err := service.DbInterface.FindArticlesByCategory(ctx, constants.NEWS, int64(articleLimit), category)
	if err != nil {
		service.Logger.Error("Failed to fetch news articles by category", "error", err)
		return nil, err
	}
	service.Logger.Info(fmt.Sprintf("Fetched %d news articles by category from database and creating summaries...", len(articles)))

	// Summarize the articles
	articles, err = service.ArticleSummaryHelper(ctx, articles)
	if err != nil {
		service.Logger.Error("Failed to summarize articles", "error", err)
		return nil, err
	}
	service.Logger.Info(fmt.Sprintf("Summarized %d news articles by category", len(articles)))

	return articles, nil
}

func (service *NewsService) ScoreNewsService(
	ctx context.Context,
	threshold float64,
	articleLimit int,
) ([]newsArticle.NewsArticleDBResponse, error) {
	service.Logger.Debug("'Service Layer': Fetching news articles by score...")

	articles, err := service.DbInterface.FindArticlesByScore(ctx, constants.NEWS, int64(articleLimit), threshold)
	if err != nil {
		service.Logger.Error("Failed to fetch news articles by score", "error", err)
		return nil, err
	}

	service.Logger.Info(fmt.Sprintf("Fetched %d news articles by score from database and creating summaries...", len(articles)))
	// Summarize the articles
	articles, err = service.ArticleSummaryHelper(ctx, articles)
	if err != nil {
		service.Logger.Error("Failed to summarize articles", "error", err)
		return nil, err
	}
	service.Logger.Info(fmt.Sprintf("Summarized %d news articles by score", len(articles)))

	return articles, nil
}

func (service *NewsService) SearchNewsService(
	ctx context.Context,
	query string,
	articleLimit int,
) ([]newsArticle.NewsArticleDBResponse, error) {
	service.Logger.Debug("'Service Layer': Searching news articles...")

	systemMessage := constants.ARTICLE_NEWS_ENTITIES_AND_INTENT_SYSTEM_PROMPT
	userMessage := fmt.Sprintf(constants.ARTICLE_NEWS_ENTITIES_AND_INTENT_USER_PROMPT, query)

	llmOutput, err := service.LLMService.ExtractEntitiesAndIntent(
		ctx,
		systemMessage,
		userMessage,
	)
	if err != nil {
		service.Logger.Error("Failed to extract entities and intent from user query", "error", err)
		return nil, err
	}
	service.Logger.Debug("Extracted entities and intent", "entities", llmOutput.Entities, "intent", llmOutput.Intent)

	articles := []newsArticle.NewsArticleDBResponse{}
	searchableQuery := strings.Join(llmOutput.Keywords, " ")
	intent := llmOutput.Intent

	switch intent {
	case "category":
		// Handle category news intent
		articles, err = service.DbInterface.FindArticlesByCategory(ctx, constants.NEWS, int64(articleLimit), llmOutput.Entities[0])
		if err != nil {
			service.Logger.Error("Failed to fetch news articles by category", "error", err)
			return nil, err
		}

	case "source":
		// Handle news by source intent
		articles, err = service.DbInterface.FindArticlesBySource(ctx, constants.NEWS, int64(articleLimit), llmOutput.Entities[0])
		if err != nil {
			service.Logger.Error("Failed to fetch news articles by source", "error", err)
			return nil, err
		}

	case "nearby":
		// Handle nearby news intent
		// find the lat and lon from the entity where the location is mentioned
		// And in case of no valid location found, fallback to normal search
		radius := 1.0 // Default radius is 1km
		latitude, longitude := 0.0, 0.0
		isFound := false
		// extract latitude and longitude from location entities
		for _, location := range llmOutput.Entities {
			lat, lon, err := utils.ExtractLatAndLon(location)
			if err == nil {
				latitude, longitude = lat, lon
				isFound = true
				break
			}
			service.Logger.Warn("Failed to extract latitude and longitude from location", "location", location, "error", err)
		}
		// If no valid location found, fallback to normal search
		if isFound == false {
			service.Logger.Error("No valid location found with respect to user query", "locations: ", llmOutput.Entities)
			service.Logger.Warn("Fallback to 'Normal Search on title and description'")
			// Fallback to normal search if no valid location found
			articles, err = service.DbInterface.FindArticlesBySearchQuery(ctx, constants.NEWS, int64(articleLimit), searchableQuery)
			if err != nil {
				service.Logger.Error("Failed to search news articles", "error", err)
				return nil, err
			}
		} else {
			// If valid location found, search for nearby articles with latitude and longitude
			articles, err = service.DbInterface.FindArticlesNearby(ctx, constants.NEWS, int64(articleLimit), latitude, longitude, radius)
			if err != nil {
				service.Logger.Error("Failed to fetch nearby news articles", "error", err)
				return nil, err
			}
		}

	case "search":
		// Handle search intent
		articles, err = service.DbInterface.FindArticlesBySearchQuery(ctx, constants.NEWS, int64(articleLimit), searchableQuery)
		if err != nil {
			service.Logger.Error("Failed to search news articles", "error", err)
			return nil, err
		}
	default:
		service.Logger.Warn("Unknown intent, Fallback to normal search", "intent", intent)
		// Fallback to normal search if intent is unknown
		articles, err = service.DbInterface.FindArticlesBySearchQuery(ctx, constants.NEWS, int64(articleLimit), searchableQuery)
		if err != nil {
			service.Logger.Error("Failed to search news articles", "error", err)
			return nil, err
		}
	}

	service.Logger.Info(fmt.Sprintf("Fetched %d news articles based on user query and creating summaries...", len(articles)))
	// Summarize the articles
	articles, err = service.ArticleSummaryHelper(ctx, articles)
	if err != nil {
		service.Logger.Error("Failed to summarize articles", "error", err)
		return nil, err
	}
	service.Logger.Info(fmt.Sprintf("Summarized %d news articles based on user query", len(articles)))

	return articles, nil
}

func (service *NewsService) SourceNewsService(
	ctx context.Context,
	source string,
	articleLimit int,
) ([]newsArticle.NewsArticleDBResponse, error) {
	service.Logger.Debug("'Service Layer': Fetching news articles by source...")

	articles, err := service.DbInterface.FindArticlesBySource(ctx, constants.NEWS, int64(articleLimit), source)
	if err != nil {
		service.Logger.Error("Failed to fetch news articles by source", "error", err)
		return nil, err
	}

	service.Logger.Info(fmt.Sprintf("Fetched %d news articles by source from database and creating summaries...", len(articles)))
	// Summarize the articles
	articles, err = service.ArticleSummaryHelper(ctx, articles)
	if err != nil {
		service.Logger.Error("Failed to summarize articles", "error", err)
		return nil, err
	}
	service.Logger.Info(fmt.Sprintf("Summarized %d news articles by source", len(articles)))

	return articles, nil
}

func (service *NewsService) NearbyNewsService(
	ctx context.Context,
	latitude float64,
	longitude float64,
	radius float64,
	articleLimit int,
) ([]newsArticle.NewsArticleDBResponse, error) {
	service.Logger.Debug("'Service Layer': Fetching nearby news articles...")

	articles, err := service.DbInterface.FindArticlesNearby(ctx, constants.NEWS, int64(articleLimit), latitude, longitude, radius)
	if err != nil {
		service.Logger.Error("Failed to fetch nearby news articles", "error", err)
		return nil, err
	}

	service.Logger.Info(fmt.Sprintf("Fetched %d nearby news articles from database and creating summaries...", len(articles)))
	// Summarize the articles
	articles, err = service.ArticleSummaryHelper(ctx, articles)
	if err != nil {
		service.Logger.Error("Failed to summarize articles", "error", err)
		return nil, err
	}
	service.Logger.Info(fmt.Sprintf("Summarized %d nearby news articles", len(articles)))

	return articles, nil
}

func (service *NewsService) SimulateEventsService(
	ctx context.Context,
	userID string,
	articleID string,
	eventType string,
	latitude float64,
	longitude float64,
) error {
	service.Logger.Debug("'Service Layer': Simulating events...")

	event := newsArticle.UserEvent{
		ArticleID: articleID,
		EventType: eventType,
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Now(),
	}

	// insert in user collection : userEvents
	err := service.DbInterface.InsertUserEvent(ctx, event)
	if err != nil {
		service.Logger.Error("Failed to insert user event", "error", err)
		return err
	}

	service.Logger.Info(fmt.Sprintf("Simulated user event for user %s", userID))
	return nil
}

func (service *NewsService) TrendingNewsService(
	ctx context.Context,
	articleLimit int,
	latitude float64,
	longitude float64,
	radius float64,
) ([]newsArticle.NewsArticleDBResponse, error) {
	service.Logger.Debug("'Service Layer': Fetching trending news articles...")

	//TODO: Caching can be implemented here to store and retrieve trending articles efficiently

	articles, err := service.DbInterface.FindTrendingArticles(ctx, constants.NEWS, constants.USER_EVENT, int64(articleLimit), latitude, longitude, radius)
	if err != nil {
		service.Logger.Error("Failed to fetch trending news articles", "error", err)
		return nil, err
	}

	service.Logger.Info(fmt.Sprintf("Fetched %d trending news articles from database and creating summaries...", len(articles)))
	// Summarize the articles
	articles, err = service.ArticleSummaryHelper(ctx, articles)
	if err != nil {
		service.Logger.Error("Failed to summarize articles", "error", err)
		return nil, err
	}
	service.Logger.Info(fmt.Sprintf("Summarized %d trending news articles", len(articles)))

	return articles, nil
}
