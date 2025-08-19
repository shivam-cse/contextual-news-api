package v1

import (
	"strconv"
	"log/slog"
	"net/http"
	"github.com/gin-gonic/gin"
	// "github.com/shivam-cse/contextual-news-api/pkg/constants"
	"github.com/shivam-cse/contextual-news-api/internal/models/newsResponse"
	"github.com/shivam-cse/contextual-news-api/internal/services"
)

type NewsHandler struct {
	NewsService   	*services.NewsService
	Logger 	  		*slog.Logger
}

func NewNewsHandler(
	newsService *services.NewsService,
	logger *slog.Logger,
) *NewsHandler {
	return &NewsHandler{
		NewsService: newsService,
		Logger:      logger,
	}
}

func (newsHandler *NewsHandler) LatestNewsHandler(c *gin.Context) {
	newsHandler.Logger.Debug("'Handler layer': Fetching latest news articles...")

	ctx := c.Request.Context()
	maxArticleLimit, err := strconv.Atoi(c.Query("articleLimit"))
	if err != nil {
		maxArticleLimit = 5 // Default to 5 if maxArticleLimit is not provided or invalid
	}

	newsArticles, err := newsHandler.NewsService.LatestNewsService(ctx, maxArticleLimit)

	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger, 
			http.StatusInternalServerError, 
			"Failed to retrieve news articles", 
			err,
		)
		return
	}

	newsResponse.Success(
		c,
		newsHandler.Logger,
		http.StatusOK,
		"Successfully retrieved all news articles",
		newsArticles,
		len(newsArticles),
	)
}

func (newsHandler *NewsHandler) CategoryNewsHandler(c *gin.Context) {
	newsHandler.Logger.Debug("'Handler layer': Fetching news articles by category...")

	ctx := c.Request.Context()
	category := c.Query("category")
	maxArticleLimit, err := strconv.Atoi(c.Query("articleLimit"))
	
	if err != nil {
		maxArticleLimit = 5 // Default to 5 if maxArticleLimit is not provided or invalid
	}

	if category == "" {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Category parameter is required",
			nil,
		)
		return
	}

	results, err := newsHandler.NewsService.CategoryNewsService(ctx, category, maxArticleLimit)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusInternalServerError,
			"Failed to retrieve news articles by category",
			err,
		)
		return
	}

	newsResponse.Success(
		c,
		newsHandler.Logger,
		http.StatusOK,
		"Successfully retrieved news articles by category",
		results,
		len(results),
	)
}

func (newsHandler *NewsHandler) ScoreNewsHandler(c *gin.Context) {
	newsHandler.Logger.Debug("'Handler layer': Fetching news articles by score...")
	ctx := c.Request.Context()
	thresholdStr := c.DefaultQuery("threshold", "0.7")
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil || threshold < 0 || threshold > 1 {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Threshold parameter must be a float between 0 and 1",
			err,
		)
		return
	}
	maxArticleLimit, err := strconv.Atoi(c.Query("articleLimit"))
	if err != nil {
		maxArticleLimit = 5 // Default to 5 if maxArticleLimit is not provided or invalid
	}

	results, err := newsHandler.NewsService.ScoreNewsService(ctx, threshold, maxArticleLimit)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusInternalServerError,
			"Failed to retrieve news articles by score",
			err,
		)
		return
	}

	newsResponse.Success(
		c,
		newsHandler.Logger,
		http.StatusOK,
		"Successfully retrieved news articles by score",
		results,
		len(results),
	)
}

func (newsHandler *NewsHandler) SearchNewsHandler(c *gin.Context) {
	newsHandler.Logger.Debug("'Handler layer': Fetching news articles by query...")
	ctx := c.Request.Context()
	query := c.Query("query")
	maxArticleLimit, err := strconv.Atoi(c.Query("articleLimit"))
	if err != nil {
		maxArticleLimit = 5 // Default to 5 if maxArticleLimit is not provided or invalid
	}

	if query == "" {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Query parameter is required",
			nil,
		)
		return
	}

	results, err := newsHandler.NewsService.SearchNewsService(ctx, query, maxArticleLimit)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusInternalServerError,
			"Failed to search news articles",
			err,
		)
		return
	}

	newsResponse.Success(
		c,
		newsHandler.Logger,
		http.StatusOK,
		"Successfully searched news articles",
		results,
		len(results),
	)
}

func (newsHandler *NewsHandler) SourceNewsHandler(c *gin.Context) {
	newsHandler.Logger.Debug("'Handler layer': Fetching news articles by source...")
	ctx := c.Request.Context()
	source := c.Query("source")
	maxArticleLimit, err := strconv.Atoi(c.Query("articleLimit"))
	if err != nil {
		maxArticleLimit = 5 // Default to 5 if maxArticleLimit is not provided or invalid
	}

	if source == "" {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Source parameter is required",
			nil,
		)
		return
	}

	results, err := newsHandler.NewsService.SourceNewsService(ctx, source, maxArticleLimit)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusInternalServerError,
			"Failed to retrieve news articles by source",
			err,
		)
		return
	}

	newsResponse.Success(
		c,
		newsHandler.Logger,
		http.StatusOK,
		"Successfully retrieved news articles by source",
		results,
		len(results),
	)
}

func (newsHandler *NewsHandler) NearbyNewsHandler(c *gin.Context) {
	newsHandler.Logger.Debug("'Handler layer': Fetching nearby news articles...")
	ctx := c.Request.Context()
	latitudeStr := c.Query("lat")
	longitudeStr := c.Query("lon")

	// Radius must be provide in kilometers
	// Default radius is 1km if not provided
	radiusStr := c.DefaultQuery("radius", "1")

	maxArticleLimit, err := strconv.Atoi(c.Query("articleLimit"))
	if err != nil {
		maxArticleLimit = 5 // Default to 5 if maxArticleLimit is not provided or invalid
	}

	if latitudeStr == "" || longitudeStr == "" {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Latitude as 'lat' and longitude as 'lon' parameters are required",
			nil,
		)
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Invalid latitude parameter",
			err,
		)
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Invalid longitude parameter",
			err,
		)
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Invalid radius parameter",
			err,
		)
		return
	}

	results, err := newsHandler.NewsService.NearbyNewsService(ctx, latitude, longitude, radius, maxArticleLimit)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusInternalServerError,
			"Failed to retrieve nearby news articles",
			err,
		)
		return
	}

	newsResponse.Success(
		c,
		newsHandler.Logger,
		http.StatusOK,
		"Successfully retrieved nearby news articles",
		results,
		len(results),
	)
}

func (newsHandler *NewsHandler) SimulateEventsHandler(c *gin.Context) {
	newsHandler.Logger.Debug("'Handler layer': Simulating events...")
	ctx := c.Request.Context()
	var req struct {
		UserID    string  `json:"user_id"`
		ArticleID string  `json:"article_id"`
		EventType string  `json:"event_type"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Invalid request payload",
			err,
		)
		return
	}

	err := newsHandler.NewsService.SimulateEventsService(ctx, req.UserID, req.ArticleID, req.EventType, req.Latitude, req.Longitude)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusInternalServerError,
			"Failed to simulate events",
			err,
		)
		return
	}

	newsResponse.Success(
		c,
		newsHandler.Logger,
		http.StatusOK,
		"Successfully simulated events",
		nil,
		0,
	)
}

func (newsHandler *NewsHandler) TrendingNewsHandler(c *gin.Context) {
	newsHandler.Logger.Debug("'Handler layer': Fetching trending news articles...")
	ctx := c.Request.Context()

	radiusStr := c.DefaultQuery("radius", "1") // Default radius is 1km if not provided
	latitudeStr := c.Query("lat")
	longitudeStr := c.Query("lon")
	maxArticleLimit, err := strconv.Atoi(c.Query("articleLimit"))

	if err != nil {
		maxArticleLimit = 5 // Default to 5 if maxArticleLimit is not provided or invalid
	}

	if latitudeStr == "" || longitudeStr == "" {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Latitude as 'lat' and longitude as 'lon' parameters are required",
			nil,
		)
		return
	}

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Invalid latitude parameter",
			err,
		)
		return
	}

	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Invalid longitude parameter",
			err,
		)
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusBadRequest,
			"Invalid radius parameter",
			err,
		)
		return
	}

	results, err := newsHandler.NewsService.TrendingNewsService(ctx, maxArticleLimit, latitude, longitude, radius)
	if err != nil {
		newsResponse.Error(
			c,
			newsHandler.Logger,
			http.StatusInternalServerError,
			"Failed to retrieve trending news articles",
			err,
		)
		return
	}

	newsResponse.Success(
		c,
		newsHandler.Logger,
		http.StatusOK,
		"Successfully retrieved trending news articles",
		results,
		len(results),
	)
}