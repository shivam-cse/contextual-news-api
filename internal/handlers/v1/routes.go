package v1

import (
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/shivam-cse/contextual-news-api/internal/models/newsResponse"
)

var DefaultTimeoutDuration = 600*time.Second

func RegisterRoutes(router *gin.Engine, newsHandlers *NewsHandler) {
	api := router.Group("/api/v1")
	{
		news := api.Group("/news")
		{
			// GET /api/v1/news/latest?articleLimit=<limit>
			news.GET("/latest", timeout.New(
				timeout.WithTimeout(DefaultTimeoutDuration),
				timeout.WithResponse(newsResponse.TimeOut),
			), newsHandlers.LatestNewsHandler)

			// GET /api/v1/news/category?category=<category>&articleLimit=<limit>
			news.GET("/category", timeout.New(
					timeout.WithTimeout(DefaultTimeoutDuration),
					timeout.WithResponse(newsResponse.TimeOut),
				), newsHandlers.CategoryNewsHandler)

			// GET /api/v1/news/category?score=<score>&articleLimit=<limit>
			news.GET("/score", timeout.New(
					timeout.WithTimeout(DefaultTimeoutDuration),
					timeout.WithResponse(newsResponse.TimeOut),
				), newsHandlers.ScoreNewsHandler)

			// GET /api/v1/news/search?query=<query>&articleLimit=<limit>
			news.GET("/search", timeout.New(
					timeout.WithTimeout(DefaultTimeoutDuration),
					timeout.WithResponse(newsResponse.TimeOut),
				), newsHandlers.SearchNewsHandler)

			// GET /api/v1/news/source?source=<source>&articleLimit=<limit>
			news.GET("/source", timeout.New(
					timeout.WithTimeout(DefaultTimeoutDuration),
					timeout.WithResponse(newsResponse.TimeOut),
				), newsHandlers.SourceNewsHandler)

			// GET /api/v1/news/nearby?lat=<latitude>&long=<longitude>&articleLimit=<limit>
			news.GET("/nearby", timeout.New(
					timeout.WithTimeout(DefaultTimeoutDuration),
					timeout.WithResponse(newsResponse.TimeOut),
				), newsHandlers.NearbyNewsHandler)

			// GET /api/v1/news/trending?lat=<latitude>&long=<longitude>&articleLimit=<limit>
			news.GET("/trending", timeout.New(
					timeout.WithTimeout(DefaultTimeoutDuration),
					timeout.WithResponse(newsResponse.TimeOut),
				), newsHandlers.TrendingNewsHandler)

			// GET /api/v1/news/trending
			news.POST("/events/simulate", timeout.New(
					timeout.WithTimeout(DefaultTimeoutDuration),
					timeout.WithResponse(newsResponse.TimeOut),
				), newsHandlers.SimulateEventsHandler)
		}
	}
}

